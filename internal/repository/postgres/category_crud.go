package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alex60217101990/test_api/internal/helpers"
	"golang.org/x/sync/errgroup"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func (r *Repository) GetCategories(ctx context.Context,
	pagination *models.Pagination, sortedBy *models.SortedBy, associate ...struct{}) (categories []*models.Category, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "get categories [pag: %+v, sort: %+v] postgres repo", pagination, sortedBy)
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
			`AND pc.id > COALESCE(
				(SELECT id FROM product_categories
				WHERE public_id = '%s'), 0
			) %s %s`, pagination.From, new(models.Category).ConvertToQuery(sortedBy),
			getCurrentLimmit(pagination.To))
	} else {
		subQueryPag = new(models.Category).ConvertToQuery(sortedBy)
	}

	var rows pgx.Rows
	rows, err = conn.Query(ctx,
		fmt.Sprintf(`
		SELECT pc.* %s 
		FROM product_categories pc
		WHERE pc.deleted_at IS NULL %s`,
			categoryAssocSubQuery(associate...), subQueryPag,
		),
	)

	if err != nil {
		return nil, err
	}

	categories = make([]*models.Category, 0)
	categories, err = repository.ConvertMultipleRowsToCategories(ctx, rows, associate...)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, pgx.ErrNoRows
	}

	return categories, err
}

func (r *Repository) GetCategoryByNameOrID(ctx context.Context, query string, associate ...struct{}) (category *models.Category, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "get category by name or id [%s] postgres repo", query)
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
		SELECT pc.* %s 
		FROM product_categories pc
		WHERE pc.deleted_at IS NULL
			AND (pc.name = '%s' OR pc.public_id = '%s')
		LIMIT 1;`,
			categoryAssocSubQuery(associate...), query, queryPublicID,
		),
	)

	if err != nil {
		return nil, err
	}

	categories := make([]*models.Category, 0)
	categories, err = repository.ConvertMultipleRowsToCategories(ctx, rows, associate...)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, pgx.ErrNoRows
	}

	return categories[0], err
}

func (r *Repository) InsertCategory(ctx context.Context, category *models.Category) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "insert new category [%+v] posrgres repo", category)
		}
	}()

	err = category.BeforeCreate()
	if err != nil {
		return err
	}

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok {
		category.ChangeByUser = uint(user.ID)
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
			`INSERT INTO product_categories (public_id, change_by_user, name, popularity) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (name) DO 
			UPDATE SET 
				change_by_user = EXCLUDED.change_by_user,  
				popularity = EXCLUDED.popularity,
				updated_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc');`,
			category.PublicID, category.ChangeByUser, category.Name, category.Popularity,
		)
		return err
	})

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			fmt.Sprintf(`WITH helper AS (
				SELECT pc.id FROM product_categories pc
				WHERE (pc.public_id = '%s' OR pc.name = '%s')
					AND pc.deleted_at IS NULL)
			INSERT INTO product_category_relations (product_id, category_id)
			SELECT pr.id, h.id
			FROM products pr, helper h
			WHERE pr.public_id IN (%s)
				AND pr.deleted_at IS NULL`,
				category.GetPublicID(), category.Name,
				repository.ConvertObjSliceToQueryStr(category.Products)),
		)
		return err
	})

	return eg.Wait()
}

func (r *Repository) UpdateCategory(ctx context.Context, category *models.Category) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "update new category [%+v] posrgres repo", category)
		}
	}()

	err = category.BeforeCreate()
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
		category.ChangeByUser = uint(user.ID)
	}

	query := fmt.Sprintf(`
	UPDATE product_categories SET
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
		category.ChangeByUser, category.Name, category.Popularity,
		category.GetPublicID(), category.Name, category.ChangeByUser,
		category.ChangeByUser, category.ChangeByUser,
		category.Name, category.Name, category.Name,
		category.Popularity, category.Popularity,
	)

	_, err = conn.Exec(ctx, query)

	return err
}

func (r *Repository) DeleteSoftCategory(ctx context.Context, query string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "soft delete category [%v] posrgres repo", query)
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
			`UPDATE product_categories SET
				deleted_at = (CURRENT_TIMESTAMP at TIME ZONE 'utc')
			WHERE public_id = $1 OR name = $2;`, query, query,
		)
		return err
	})

	eg.Go(func() error {
		_, err := tx.ExecContext(ctx,
			`DELETE FROM product_category_relations WHERE category_id IN (
				SELECT pc.id FROM product_categories pc
				WHERE (pc.public_id = $1 OR pc.name = $2) 
			)`, query, query,
		)
		return err
	})

	return eg.Wait()
}

func (r *Repository) DeleteHardCategory(ctx context.Context, query string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "hard delete category [%v] posrgres repo", query)
		}
	}()

	var conn *pgxpool.Conn
	conn, err = r.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		`DELETE FROM product_categories WHERE public_id = $1 OR name = $2;`, query, query,
	)

	return err
}
