package mock

import (
	"context"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/brianvoe/gofakeit"
	"github.com/pkg/errors"
)

type Repository struct{}

func (r *Repository) Connect(ctx context.Context) error {
	gofakeit.Seed(0)
	return nil
}

func (r *Repository) Close() error {
	return nil
}

func (r *Repository) GetDB() (interface{}, error) {
	return nil, nil
}

func (r *Repository) GetStats() map[string]interface{} {
	return nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}

func (r *Repository) GetUserByCreeds(ctx context.Context, creeds *models.Credentials, associate ...struct{}) (user *models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "get user by creeds mock repo")
		}
	}()

	return r.GenerateUser(ctx)
}

func (r *Repository) GetUserByPubliID(ctx context.Context, publicID string, associate ...struct{}) (user *models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "get user by 'public_id' mock repo")
		}
	}()

	return r.GenerateUser(ctx)
}

func (r *Repository) InsertUser(ctx context.Context, user *models.User) (err error) {
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) (err error) {
	return nil
}

func (r *Repository) DeleteSoftUser(ctx context.Context, publicID string) (err error) {
	return nil
}

func (r *Repository) DeleteHardUser(ctx context.Context, publicID string) (err error) {
	return nil
}

func (r *Repository) GetCategories(ctx context.Context, pagination *models.Pagination,
	sortedBy *models.SortedBy, associate ...struct{}) (categories []*models.Category, err error) {
	return r.GenerateCategories(ctx, 0)
}

func (r *Repository) GetCategoryByNameOrID(ctx context.Context, query string, associate ...struct{}) (category *models.Category, err error) {
	return r.GenerateCategory(ctx)
}

func (r *Repository) InsertCategory(ctx context.Context, category *models.Category) (err error) {
	return nil
}

func (r *Repository) UpdateCategory(ctx context.Context, category *models.Category) (err error) {
	return nil
}

func (r *Repository) DeleteSoftCategory(ctx context.Context, query string) (err error) {
	return nil
}

func (r *Repository) DeleteHardCategory(ctx context.Context, query string) (err error) {
	return nil
}

func (r *Repository) GetProducts(ctx context.Context,
	pagination *models.Pagination, sortedBy *models.SortedBy,
	associate ...struct{}) (products []*models.Product, err error) {
	return r.GenerateProducts(ctx, 0)
}

func (r *Repository) GetProductByNameOrID(ctx context.Context, query string, associate ...struct{}) (product *models.Product, err error) {
	return r.GenerateProduct(ctx)
}

func (r *Repository) InsertProduct(ctx context.Context, product *models.Product) (err error) {
	return nil
}
func (r *Repository) UpdateProduct(ctx context.Context, product *models.Product) (err error) {
	return nil
}

func (r *Repository) DeleteSoftProduct(ctx context.Context, query string) (err error) {
	return nil
}
func (r *Repository) DeleteHardProduct(ctx context.Context, query string) (err error) {
	return nil
}

func (r *Repository) AddRelationCategory(ctx context.Context, productID, catgoryID string) (err error) {
	return nil
}
func (r *Repository) DelRelationCategory(ctx context.Context, productID, catgoryID string) (err error) {
	return nil
}
