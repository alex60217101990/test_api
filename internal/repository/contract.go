package repository

import (
	"context"

	"github.com/alex60217101990/test_api/internal/models"
)

type SQLRow interface {
	Scan(dest ...interface{}) error
}

type SQLRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type Repository interface {
	Connect(ctx context.Context) error
	Close()
	GetDB() (interface{}, error)
	GetStats() map[string]interface{}
	Ping(ctx context.Context) error
	// user methods
	GetUserByCreeds(ctx context.Context, creeds *models.Credentials, associate ...struct{}) (user *models.User, err error)
	InsertUser(ctx context.Context, user *models.User) (err error)
	UpdateUser(ctx context.Context, user *models.User) (err error)
	DeleteSoft(ctx context.Context, publicID string) (err error)
	DeleteHard(ctx context.Context, publicID string) (err error)
}
