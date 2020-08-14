package mock

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

func (r *Repository) GenerateCreeds(ctx context.Context) *models.Credentials {
	return &models.Credentials{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 32),
	}
}

func (r *Repository) GenerateUser(ctx context.Context) (u *models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new user struct")
		}
	}()

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	oldDate := TimeToTimePtr(gofakeit.DateRange(time.Now().AddDate(0, -5, 0), time.Now().AddDate(0, -3, 0)))

	return &models.User{
		Base: models.Base{
			ID:        gofakeit.Int64(),
			PublicID:  uuid,
			CreatedAt: oldDate,
			UpdatedAt: oldDate,
		},
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 32),
		IsOnline: gofakeit.Bool(),
	}, nil
}
