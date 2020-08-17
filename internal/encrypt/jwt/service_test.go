package jwt

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/pkg/errors"

	"github.com/alex60217101990/test_api/internal/encrypt"
	"github.com/alex60217101990/test_api/internal/helpers"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/repository/mock"
)

var (
	service  encrypt.SecretService
	confFile = "../../../deploy/configs/application.yaml"
)

func init() {
	ctx := context.Background()
	// Load configs file
	err := configs.ReadConfigFile(confFile)
	if err != nil {
		logger.CmdError.Println(err)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	// Init loggers
	logger.InitLoggerSettings()

	m := &mock.Repository{}
	m.Connect(ctx)

	configs.Conf.Keys.PubKeyAuth = ".." + string(os.PathSeparator) + configs.Conf.Keys.PubKeyAuth
	configs.Conf.Keys.PrvKeyAuth = ".." + string(os.PathSeparator) + configs.Conf.Keys.PrvKeyAuth

	service = NewSecretService(ctx, m)
}

func TestGenerateToken(t *testing.T) {
	user := &models.User{
		Base: models.Base{
			ID: 1,
		},
	}
	user.FromStr("6e102015-03a0-4c77-b2bc-9ab73eb3773a")
	ctx := helpers.AddToContext(context.Background(), repository.UserSessionKey, user)
	t.Run("success", func(t *testing.T) {
		token, err := service.GenerateToken(ctx)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(token)
		t.Log(token)
	})

	t.Run("failed", func(t *testing.T) {
		configs.Conf.APIKey = " "
		token, err := service.GenerateToken(ctx)
		if err != nil && !assert.Contains(t, err.Error(), "generate JWT token") {
			t.Error(err)
			return
		}
		t.Log(token)
	})
}

func TestVerifyTokenString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := service.VerifyTokenString(context.Background(), `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI2ZTEwMjAxNS0wM2EwLTRjNzctYjJiYy05YWI3M2ViMzc3M2EiLCJleHAiOjE1OTc3MzQ5ODEsInN1YiI6IjM0N3k3NTQ3NTg3N040N2MreTQ3cjR0N2M4Q3kzNDczYjQ3MjQ2NTYyNTN2MnI1MXI1MkAjXHUwMDNlXG4ifQ.bDPzOQaNxovWGioqNuaB0bZu_1fGAS3KcZ5DRw2XzBE`)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("failed", func(t *testing.T) {
		_, err := service.VerifyTokenString(context.Background(), "some_token")
		if err != nil && !assert.Contains(t, err.Error(), "verify JWT token") {
			t.Error(err)
			return
		}
	})
}

func TestRefreshToken(t *testing.T) {
	token, err := service.RefreshToken(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI2ZTEwMjAxNS0wM2EwLTRjNzctYjJiYy05YWI3M2ViMzc3M2EiLCJleHAiOjE1OTc3MzQ5ODEsInN1YiI6IjM0N3k3NTQ3NTg3N040N2MreTQ3cjR0N2M4Q3kzNDczYjQ3MjQ2NTYyNTN2MnI1MXI1MkAjXHUwMDNlXG4ifQ.bDPzOQaNxovWGioqNuaB0bZu_1fGAS3KcZ5DRw2XzBE")
	if err != nil {
		if assert.NotEqual(t, errors.Cause(err), jwt.ErrSignatureInvalid) {
			return
		}
		t.Error(err)
	}
	t.Log(token)
}
