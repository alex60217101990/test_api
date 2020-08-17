package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/alex60217101990/test_api/internal/helpers"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func (r *Repository) GetProducts(ctx context.Context,
	pagination *models.Pagination, sortedBy *models.SortedBy, associate ...struct{}) (products []*models.Product, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "get products [pag: %+v, sort: %+v] postgres repo", pagination, sortedBy)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var subQueryPag string

	if pagination != nil {
		pagination.From, err = getCorrectUUID(pagination.From)
		if err != nil {
			return nil, err
		}

		subQueryPag = fmt.Sprintf(
			`AND pr.id > COALESCE(
				(SELECT id FROM products
				WHERE public_id = '%s' OR name = '%s'), 0
			) %s %s`, pagination.From, pagination.From,
			new(models.Product).ConvertToQuery(sortedBy),
			getCurrentLimmit(pagination.To))
	} else {
		subQueryPag = new(models.Product).ConvertToQuery(sortedBy)
	}

	var rows pgx.Rows
	rows, err = conn.Query(ctx,
		fmt.Sprintf(`
		SELECT pr.* %s 
		FROM products pr
		WHERE pr.deleted_at IS NULL %s`,
			productAssocSubQuery(associate...), subQueryPag,
		),
	)

	if err != nil {
		return nil, err
	}

	products = make([]*models.Product, 0)
	products, err = repository.ConvertMultipleRowsToProducts(ctx, rows, associate...)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, pgx.ErrNoRows
	}

	return products, err
}

func (r *Repository) GetProductByNameOrID(ctx context.Context, query string, associate ...struct{}) (product *models.Product, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "get product by name or id [%s] postgres repo", query)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var queryPublicID string
	queryPublicID, err = getCorrectUUID(query)
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows
	rows, err = conn.Query(ctx,
		fmt.Sprintf(`
		SELECT pr.* %s 
		FROM products pr
		WHERE pr.deleted_at IS NULL
			AND (pr.name = '%s' OR pr.public_id = '%s')`,
			productAssocSubQuery(associate...), query, queryPublicID,
		),
	)

	if err != nil {
		return nil, err
	}

	products := make([]*models.Product, 0)
	products, err = repository.ConvertMultipleRowsToProducts(ctx, rows, associate...)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, pgx.ErrNoRows
	}

	return products[0], err
}

func (r *Repository) InsertProduct(ctx context.Context, product *models.Product) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "insert new product [%+v] posrgres repo", product)
		}
	}()

	err = product.BeforeCreate()
	if err != nil {
		return err
	}

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok {
		product.ChangeByUser = uint(user.ID)
	}

	var tx *sql.Tx
	tx, err = r.initTx(ctx)
	defer r.closeTx(ctx, tx, err)
	if err != nil {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			fmt.Sprintf(`WITH helper AS (
				SELECT pr.id FROM products pr
				WHERE (pr.public_id = '%s' OR pr.name = '%s')
					AND pr.deleted_at IS NULL)
			INSERT INTO product_category_relations (product_id, category_id) 
			SELECT h.id, pc.id 
			FROM product_categories pc, helper h
			WHERE pc.public_id IN (%s) 
				AND pc.deleted_at IS NULL`,
				product.GetPublicID(), product.Name,
				repository.ConvertObjSliceToQueryStr(product.Categories)),
		)
		return err
	})

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO products (public_id, change_by_user, name, popularity) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (name) DO 
			UPDATE SET 
				change_by_user = EXCLUDED.change_by_user,  
				popularity = EXCLUDED.popularity,
				updated_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc');`,
			product.PublicID, product.ChangeByUser, product.Name, product.Popularity,
		)
		return err
	})

	return eg.Wait()
}

func (r *Repository) UpdateProduct(ctx context.Context, product *models.Product) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "update product [%+v] posrgres repo", product)
		}
	}()

	err = product.BeforeCreate()
	if err != nil {
		return err
	}

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok {
		product.ChangeByUser = uint(user.ID)
	}

	query := fmt.Sprintf(`
	UPDATE products SET
		change_by_user = COALESCE(%d, change_by_user),
		name = COALESCE('%s', name),
		popularity = COALESCE(%d, popularity),
		updated_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc')
	WHERE (public_id = '%s' OR name = '%s')
		AND (%d IS NOT NULL AND %d IS DISTINCT FROM change_by_user AND (%d > 0)
			OR '%s' IS NOT NULL AND '%s' IS DISTINCT FROM name AND (('%s' = '') IS FALSE)
			OR %d IS NOT NULL AND %d IS DISTINCT FROM popularity
		) 
		AND deleted_at IS NULL;`,
		product.ChangeByUser, product.Name, product.Popularity,
		product.GetPublicID(), product.Name, product.ChangeByUser,
		product.ChangeByUser, product.ChangeByUser,
		product.Name, product.Name, product.Name,
		product.Popularity, product.Popularity,
	)

	_, err = conn.Exec(ctx, query)

	return err
}

func (r *Repository) AddRelationCategory(ctx context.Context, productID, catgoryID string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "add category realion [%s] to product [%s] posrgres repo", catgoryID, productID)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`INSERT INTO product_category_relations (product_id, category_id) 
		VALUES ((
			SELECT pr.id FROM products pr
			WHERE pr.public_id = $1
				AND pr.deleted_at IS NULL 
		), (
			SELECT pc.id FROM product_categories pc
			WHERE pc.public_id = $2
				AND pc.deleted_at IS NULL
		))`, productID, catgoryID,
	)

	return err
}

func (r *Repository) DelRelationCategory(ctx context.Context, productID, catgoryID string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "del category realion [%s] from product [%s] posrgres repo", productID, catgoryID)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		fmt.Sprintf(`DELETE FROM product_category_relations WHERE product_id = (
			SELECT pr.id FROM products pr
			WHERE (pr.public_id = '%s' OR pr.name = '%s')
				AND pr.deleted_at IS NULL 
		) AND category_id = (
			SELECT pc.id FROM product_categories pc
			WHERE (pc.public_id = '%s' OR pc.name = '%s')
				AND pc.deleted_at IS NULL
		)`, productID, productID, catgoryID, catgoryID),
	)

	return err
}

func (r *Repository) DeleteSoftProduct(ctx context.Context, query string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "soft delete product [%v] posrgres repo", query)
		}
	}()

	var tx *sql.Tx
	tx, err = r.initTx(ctx)
	defer r.closeTx(ctx, tx, err)
	if err != nil {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			`UPDATE products SET
				deleted_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc')
			WHERE public_id = $1 OR name = $2;`, query, query,
		)
		return err
	})

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			`DELETE FROM product_category_relations WHERE product_id IN (
				SELECT pr.id FROM products pr
				WHERE (pr.public_id = $1 OR pr.name = $2) 
			)`, query, query,
		)
		return err
	})

	return eg.Wait()
}

func (r *Repository) DeleteHardProduct(ctx context.Context, query string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "hard delete product [%v] posrgres repo", query)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`DELETE FROM products WHERE public_id = $1 OR name = $2;`, query, query,
	)

	return err
}
