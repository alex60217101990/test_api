package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func convertMapConfigsToUserSlice(ctx context.Context, rowMap map[string]*models.User) (users []*models.User) {
	users = make([]*models.User, 0)

	for _, user := range rowMap {
		users = append(users, user)
	}

	return users
}

func convertMultipleSQLRowsToUser(ctx context.Context, rows SQLRows) (users []*models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "convert multiple SQL rows to model")
		}
	}()

	tmpUsersMap := make(map[string]*models.User)

	for rows.Next() {
		var (
			newUser     models.User
			newCategory models.Category
			newProduct  models.Product

			changeUserIDCat, changeUserIDProd, categoryID int
		)

		err = rows.Scan(
			&newUser.ID, &newUser.PublicID, &newUser.Username,
			&newUser.Email, &newUser.Password, &newUser.IsOnline,
			&newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt,
			&newCategory.ID, &newCategory.PublicID, &changeUserIDCat,
			&newCategory.Name, &newCategory.Popularity,
			&newCategory.CreatedAt, &newCategory.UpdatedAt,
			&newCategory.DeletedAt, &newProduct.ID, &newProduct.PublicID,
			&changeUserIDProd, &categoryID, &newProduct.Name,
			&newProduct.Popularity, &newProduct.CreatedAt,
			&newProduct.UpdatedAt, &newProduct.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if u, ok := tmpUsersMap[newUser.PublicID.String()]; !ok {
			if newCategory.ID > 0 {
				newUser.Categories = []*models.Category{&newCategory}
			}
			tmpUsersMap[newUser.PublicID.String()] = &newUser
		} else {
			if u.Products == nil {
				u.Products = make([]*models.Product, 0)
			}
			if newProduct.ID > 0 {
				u.Products = append(u.Products, &newProduct)
			}
			if u.Categories == nil {
				u.Categories = make([]*models.Category, 0)
			}
			if newCategory.ID > 0 {
				u.Categories = append(u.Categories, &newCategory)
			}
		}
	}

	return convertMapConfigsToUserSlice(ctx, tmpUsersMap), nil
}

func ConvertMultipleRowsToUser(ctx context.Context, row interface{}) (users []*models.User, err error) {
	users = make([]*models.User, 0)

	defer func() {
		for _, user := range users {
			if len(user.Password) > 0 {
				err = user.AfterFind()
			}
		}
		if err != nil {
			err = errors.WithMessage(err, "convert multiple rows to User/s model")
		}
	}()

	switch v := row.(type) {
	case *sql.Rows:
		return convertMultipleSQLRowsToUser(ctx, v)
	case pgx.Rows:
		return convertMultipleSQLRowsToUser(ctx, v)
	default:
		return nil, fmt.Errorf("has't info about type %T", v)
	}
}
