package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/liamylian/jsontime"
	"github.com/pkg/errors"
)

func checkEmptyJson(data []byte) error {
	err := json.Unmarshal(data, &struct{}{})
	if err != nil {
		return err
	}
	return nil
}

func convertMultipleSQLRowsToCategories(ctx context.Context, rows SQLRows, assoc ...struct{}) (categories []*models.Category, err error) {
	categories = make([]*models.Category, 0)

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "convert multiple SQL rows to Category model")
		}
	}()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var (
		userRelationRow, productsListRow string
	)
	for rows.Next() {
		newCategory := models.Category{
			Products: make([]*models.Product, 0),
		}

		if len(assoc) > 0 {
			err = rows.Scan(
				&newCategory.ID, &newCategory.PublicID, &newCategory.ChangeByUser,
				&newCategory.Name, &newCategory.Popularity,
				&newCategory.CreatedAt, &newCategory.UpdatedAt, &newCategory.DeletedAt,
				&userRelationRow, &productsListRow,
			)
		} else {
			err = rows.Scan(
				&newCategory.ID, &newCategory.PublicID, &newCategory.ChangeByUser,
				&newCategory.Name, &newCategory.Popularity,
				&newCategory.CreatedAt, &newCategory.UpdatedAt, &newCategory.DeletedAt,
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
				newCategory.User = &user
			}
		}

		if len(productsListRow) > 0 {
			err = json.Unmarshal([]byte(productsListRow), &newCategory.Products)

			if err1 := checkEmptyJson([]byte(productsListRow)); err != nil && err1 != nil {
				return nil, err
			}
		}

		categories = append(categories, &newCategory)
	}

	return categories, nil
}

func ConvertMultipleRowsToCategories(ctx context.Context, row interface{}, assoc ...struct{}) (categories []*models.Category, err error) {
	defer func() {
		for _, cat := range categories {
			if cat.User != nil {
				err = cat.User.AfterFind()
			}
		}
		if err != nil {
			err = errors.WithMessage(err, "convert multiple rows to Categor[y/ies] model")
		}
	}()

	switch v := row.(type) {
	case *sql.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToCategories(ctx, v, assoc...)
	case pgx.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToCategories(ctx, v, assoc...)
	default:
		return nil, fmt.Errorf("has't info about type %T", v)
	}
}
