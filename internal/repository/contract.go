package repository

import (
	"context"

	"github.com/alex60217101990/test_api/internal/models"
)

const (
	UserSessionKey = "user_session"
)

type SQLRow interface {
	Scan(dest ...interface{}) error
}

type HasRelations interface {
	GetPublicID() string
}

type SQLRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

type UserRepo interface {
	GetUserByCreeds(ctx context.Context, creeds *models.Credentials, associate ...struct{}) (user *models.User, err error)
	GetUserByPubliID(ctx context.Context, publicID string, associate ...struct{}) (user *models.User, err error)
	InsertUser(ctx context.Context, user *models.User) (err error)
	UpdateUser(ctx context.Context, user *models.User) (err error)
	DeleteSoftUser(ctx context.Context, publicID string) (err error)
	DeleteHardUser(ctx context.Context, publicID string) (err error)
}

type Repository interface {
	Connect(ctx context.Context) error
	Close() error
	GetDB() (interface{}, error)
	GetStats() map[string]interface{}
	Ping(ctx context.Context) error
	// user methods
	GetUserByCreeds(ctx context.Context, creeds *models.Credentials, associate ...struct{}) (user *models.User, err error)
	GetUserByPubliID(ctx context.Context, publicID string, associate ...struct{}) (user *models.User, err error)
	InsertUser(ctx context.Context, user *models.User) (err error)
	UpdateUser(ctx context.Context, user *models.User) (err error)
	DeleteSoftUser(ctx context.Context, publicID string) (err error)
	DeleteHardUser(ctx context.Context, publicID string) (err error)
	// category methods
	GetCategories(ctx context.Context,
		pagination *models.Pagination,
		sortedBy *models.SortedBy,
		associate ...struct{}) (categories []*models.Category, err error)
	GetCategoryByNameOrID(ctx context.Context, query string, associate ...struct{}) (category *models.Category, err error)
	InsertCategory(ctx context.Context, category *models.Category) (err error)
	UpdateCategory(ctx context.Context, category *models.Category) (err error)
	DeleteSoftCategory(ctx context.Context, query string) (err error)
	DeleteHardCategory(ctx context.Context, query string) (err error)
	// product methods
	GetProducts(ctx context.Context,
		pagination *models.Pagination, sortedBy *models.SortedBy,
		associate ...struct{}) (products []*models.Product, err error)
	GetProductByNameOrID(ctx context.Context, query string, associate ...struct{}) (product *models.Product, err error)
	InsertProduct(ctx context.Context, product *models.Product) (err error)
	UpdateProduct(ctx context.Context, product *models.Product) (err error)
	DeleteSoftProduct(ctx context.Context, query string) (err error)
	DeleteHardProduct(ctx context.Context, query string) (err error)
	AddRelationCategory(ctx context.Context, productID, catgoryID string) (err error)
	DelRelationCategory(ctx context.Context, productID, catgoryID string) (err error)
}
