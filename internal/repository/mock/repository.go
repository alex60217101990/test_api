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

func (r *Repository) Close() {}

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

func (r *Repository) InsertUser(ctx context.Context, user *models.User) (err error) {
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) (err error) {
	return nil
}

func (r *Repository) DeleteSoft(ctx context.Context, publicID string) (err error) {
	return nil
}

func (r *Repository) DeleteHard(ctx context.Context, publicID string) (err error) {
	return nil
}
