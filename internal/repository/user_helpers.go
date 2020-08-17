package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/liamylian/jsontime"
	"github.com/pkg/errors"

	"github.com/alex60217101990/test_api/internal/models"

	"github.com/jackc/pgx/v4"
)

func GenerateTestCredentials() []*models.Credentials {
	return []*models.Credentials{
		&models.Credentials{
			Username: "Alex",
			Password: "3733366636643635356637303631373337333737366637323634",
		},
		// &models.Credentials{
		// 	Username: "Dima",
		// 	Password: "@#dfghejer8768f9rg",
		// },
		// &models.Credentials{
		// 	Username: "contact@ibm.com",
		// 	Password: "@#dfghejer8768f9rg",
		// },
		// &models.Credentials{
		// 	Username: "Dima1",
		// 	Password: "@#dfghejer8768f9rg",
		// },
	}
}

func convertSingleSQLRowToUser(ctx context.Context, row SQLRow) (u *models.User, err error) {
	u = &models.User{}

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "convert single SQL row to User model")
		}
	}()

	err = row.Scan(&u.ID, &u.PublicID, &u.Username,
		&u.Email, &u.Password, &u.IsOnline,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ConvertSingleRowToUser(ctx context.Context, row interface{}) (u *models.User, err error) {
	u = &models.User{}

	defer func() {
		if len(u.Password) > 0 {
			err = u.AfterFind()
		}
		if err != nil {
			err = errors.WithMessage(err, "convert row to User model")
		}
	}()

	switch v := row.(type) {
	case *sql.Row:
		return convertSingleSQLRowToUser(ctx, v)
	case pgx.Row:
		return convertSingleSQLRowToUser(ctx, v)
	default:
		return nil, fmt.Errorf("has't info about type %T", v)
	}
}

func convertMultipleSQLRowsToUser(ctx context.Context, rows SQLRows) (users []*models.User, err error) {
	users = make([]*models.User, 0)

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "convert multiple SQL rows to User model")
		}
	}()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var (
		categoriesListRow, productsListRow string
	)
	for rows.Next() {
		newUser := models.User{
			Categories: make([]*models.Category, 0),
			Products:   make([]*models.Product, 0),
		}

		err = rows.Scan(
			&newUser.ID, &newUser.PublicID, &newUser.Username,
			&newUser.Email, &newUser.Password, &newUser.IsOnline,
			&newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt,
			&categoriesListRow, &productsListRow,
		)
		if err != nil {
			return nil, err
		}

		var json = jsontime.ConfigWithCustomTimeFormat

		if len(categoriesListRow) > 0 {
			err = json.Unmarshal([]byte(categoriesListRow), &newUser.Categories)

			if err1 := checkEmptyJson([]byte(categoriesListRow)); err != nil && err1 != nil {
				return nil, err
			}
		}

		if len(productsListRow) > 0 {
			err = json.Unmarshal([]byte(productsListRow), &newUser.Products)

			if err1 := checkEmptyJson([]byte(productsListRow)); err != nil && err1 != nil {
				return nil, err
			}
		}

		users = append(users, &newUser)
	}

	return users, nil
}

func afterUsers(users []*models.User) (err error) {
	for _, user := range users {
		if len(user.Password) > 0 {
			err = user.AfterFind()
		}
	}
	return err
}

func ConvertMultipleRowsToUser(ctx context.Context, row interface{}) (users []*models.User, err error) {
	users = make([]*models.User, 0)

	defer func() {
		err = afterUsers(users)
		if err != nil {
			err = errors.WithMessage(err, "convert multiple rows to User/s model")
		}
	}()

	switch v := row.(type) {
	case *sql.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToUser(ctx, v)
	case pgx.Rows:
		defer v.Close()
		return convertMultipleSQLRowsToUser(ctx, v)
	default:
		return nil, fmt.Errorf("has't info about type %T", v)
	}
}
