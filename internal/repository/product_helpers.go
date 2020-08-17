package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alex60217101990/test_api/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/liamylian/jsontime"
	"github.com/pkg/errors"
)

func convertMultipleSQLRowsToProducts(ctx context.Context, rows SQLRows, assoc ...struct{}) (products []*models.Product, err error) {
	products = make([]*models.Product, 0)

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "convert multiple SQL rows to Product model")
		}
	}()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var (
		userRelationRow, categoriesListRow string
	)
	for rows.Next() {
		newProduct := models.Product{
			Categories: make([]*models.Category, 0),
		}
		if len(assoc) > 0 {
			err = rows.Scan(
				&newProduct.ID, &newProduct.PublicID, &newProduct.ChangeByUser,
				&newProduct.Name, &newProduct.Popularity,
				&newProduct.CreatedAt, &newProduct.UpdatedAt, &newProduct.DeletedAt,
				&userRelationRow, &categoriesListRow,
			)
		} else {
			err = rows.Scan(
				&newProduct.ID, &newProduct.PublicID, &newProduct.ChangeByUser,
				&newProduct.Name, &newProduct.Popularity,
				&newProduct.CreatedAt, &newProduct.UpdatedAt, &newProduct.DeletedAt,
			)
		}
		if err != nil {
			return nil, err
		}

		var json = jsontime.ConfigWithCustomTimeFormat

		if len(userRelationRow) > 0 {
			var user models.User
			err = json.Unmarshal([]byte(userRelationRow), &user)

			if err1 := checkEmptyJson([]byte(userRelationRow)); err != nil && err1 != nil {
				return nil, err
			}

			if user.ID > 0 {
				newProduct.User = &user
			}
		}

		if len(categoriesListRow) > 0 {
			err = json.Unmarshal([]byte(categoriesListRow), &newProduct.Categories)

			if err1 := checkEmptyJson([]byte(categoriesListRow)); err != nil && err1 != nil {
				return nil, err
			}
		}

		products = append(products, &newProduct)
	}

	return products, nil
}

func ConvertMultipleRowsToProducts(ctx context.Context, row interface{}, assoc ...struct{}) (products []*models.Product, err error) {
	defer func() {
		for _, pr := range products {
			if pr.User != nil {
				err = pr.User.AfterFind()
			}
		}
		if err != nil {
			err = errors.WithMessage(err, "convert multiple rows to Products slice")
		}
	}()

	switch v := row.(type) {
	case *sql.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToProducts(ctx, v, assoc...)
	case pgx.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToProducts(ctx, v, assoc...)
	default:
		return nil, fmt.Errorf("has't info about type %T", v)
	}
}
