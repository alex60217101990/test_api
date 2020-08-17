package mock

import (
	"context"
	"time"

	"github.com/alex60217101990/test_api/internal/helpers"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Repository) GenerateCreeds(ctx context.Context) *models.Credentials {
	return &models.Credentials{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 32),
	}
}

func (r *Repository) GenerateBase(ctx context.Context) (b models.Base, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new base struct")
		}
	}()

	uuid, err := uuid.NewUUID()
	if err != nil {
		return b, err
	}

	oldDate := TimeToTimePtr(gofakeit.DateRange(time.Now().AddDate(0, -5, 0), time.Now().AddDate(0, -3, 0)))

	return models.Base{
		ID:        gofakeit.Int64(),
		PublicID:  uuid,
		CreatedAt: oldDate,
		UpdatedAt: oldDate,
	}, nil
}

func (r *Repository) GenerateCategory(ctx context.Context) (c *models.Category, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new category struct")
		}
	}()

	c = &models.Category{
		Name:       gofakeit.Word(),
		Popularity: int64(gofakeit.Number(0, 100)),
	}

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok {
		c.User = user
		c.ChangeByUser = uint(user.ID)
	}

	c.Base, err = r.GenerateBase(ctx)

	return c, err
}

func (r *Repository) GenerateCategories(ctx context.Context, number uint8) (c []*models.Category, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new categories struct slice")
		}
	}()

	if number == 0 {
		number = uint8(gofakeit.Number(1, 15))
	}

	c = make([]*models.Category, number, number)

	for i := 0; i < int(number); i++ {
		var newCat *models.Category
		newCat, err = r.GenerateCategory(ctx)
		if err != nil {
			return nil, err
		}
		c = append(c, newCat)
	}

	return c, err
}

func (r *Repository) GenerateUser(ctx context.Context) (u *models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new user struct")
		}
	}()

	u = &models.User{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 32),
		IsOnline: gofakeit.Bool(),
	}

	u.Base, err = r.GenerateBase(ctx)

	return u, err
}

func (r *Repository) GenerateUsers(ctx context.Context, number uint8) (u []*models.User, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new users struct slice")
		}
	}()

	if number == 0 {
		number = uint8(gofakeit.Number(1, 15))
	}

	u = make([]*models.User, number, number)

	for i := 0; i < int(number); i++ {
		var newUser *models.User
		newUser, err = r.GenerateUser(ctx)
		if err != nil {
			return nil, err
		}
		u = append(u, newUser)
	}

	return u, err
}

func (r *Repository) GenerateProduct(ctx context.Context) (p *models.Product, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new product struct")
		}
	}()

	p = &models.Product{
		Name:       gofakeit.Word(),
		Popularity: uint32(gofakeit.Number(0, 100)),
	}

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok {
		p.User = user
		p.ChangeByUser = uint(user.ID)
	}

	p.Base, err = r.GenerateBase(ctx)

	return p, err
}

func (r *Repository) GenerateProducts(ctx context.Context, number uint8) (p []*models.Product, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate new product struct slice")
		}
	}()

	if number == 0 {
		number = uint8(gofakeit.Number(1, 15))
	}

	p = make([]*models.Product, number, number)

	for i := 0; i < int(number); i++ {
		var newProd *models.Product
		newProd, err = r.GenerateProduct(ctx)
		if err != nil {
			return nil, err
		}
		p = append(p, newProd)
	}

	return p, err
}
