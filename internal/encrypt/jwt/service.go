package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/encrypt"
	"github.com/alex60217101990/test_api/internal/helpers"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JwtSecret struct {
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	Token     *jwt.Token
	repo      repository.UserRepo
}

func NewSecretService(ctx context.Context, repo repository.Repository) encrypt.SecretService {
	service := &JwtSecret{
		repo: repo,
	}
	err := service.Init(ctx)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
	return service
}

func (s *JwtSecret) Init(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "init JWT secret service")
		}
	}()

	s.signKey, s.verifyKey, err = encrypt.GetRSAKeys(ctx,
		configs.Conf.Keys.PubKeyAuth,
		configs.Conf.Keys.PrvKeyAuth,
	)
	if err != nil {
		currentDir, err := os.Getwd()
		if err != nil {
			logger.AppLogger.Fatal(err)
		}

		currentDir = ".." + currentDir
		err = encrypt.InitKeys("test-api-auth.rsa", currentDir, true, 4<<10)
		if err != nil {
			return err
		}
		s.signKey, s.verifyKey, err = encrypt.GetRSAKeys(ctx,
			configs.Conf.Keys.PubKeyAuth,
			configs.Conf.Keys.PrvKeyAuth,
		)
	}

	logger.CmdInfo.Println("🔐 Secret keys load (init) success.")

	return err
}

func (s *JwtSecret) GenerateToken(ctx context.Context) (token string, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "generate JWT token")
		}
	}()

	// set claims
	claims := jwt.StandardClaims{}
	// decode API key
	claims.Subject, err = encrypt.DecodeAPIKey()
	if err != nil {
		return token, err
	}

	data := helpers.GetFromContext(ctx, repository.UserSessionKey)
	if user, ok := data.(*models.User); ok && user != nil && (!user.IsEmpty()) {
		claims.Audience = user.GetPublicID()
	} else {
		return token, fmt.Errorf("empty User model in context")
	}

	// set the expire time
	claims.ExpiresAt = time.Now().Add(helpers.TimeoutHour(configs.Conf.Timeouts.ExpHours)).Unix()

	// create a signer for rsa
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set header
	t.Header["alg"] = jwt.SigningMethodHS256.Alg()

	tokenString, err := t.SignedString(encrypt.EncodePublicKeyToPEM(s.verifyKey))
	if err != nil {
		return token, err
	}

	return tokenString, err
}

func (s *JwtSecret) parseToken(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "parse JWT token [%s]", tokenString)
		}
	}()

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return encrypt.EncodePublicKeyToPEM(s.verifyKey), nil
	})
}

func (s *JwtSecret) VerifyTokenString(ctx context.Context, tokenString string) (newCtx context.Context, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "verify JWT token [%s]", tokenString)
		}
	}()

	var token *jwt.Token
	token, err = s.parseToken(ctx, tokenString)
	if err != nil {
		return ctx, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if token.Valid {
			var apiKey string
			apiKey, err = encrypt.DecodeAPIKey()
			if err != nil {
				return ctx, err
			}

			if sub, ok := claims["sub"].(string); ok && sub == apiKey {
				if aud, ok := claims["aud"].(string); ok && len(aud) > 0 {
					user, err := s.repo.GetUserByPubliID(ctx, aud)
					if err != nil {
						return ctx, err
					}
					ctx = helpers.AddToContext(ctx, repository.UserSessionKey, user)
					return ctx, nil
				}
			}
		}
	}
	return ctx, jwt.ErrSignatureInvalid
}

func (s *JwtSecret) RefreshToken(ctx context.Context, tokenString string) (tokenStr string, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "refresh JWT token [%s]", tokenString)
		}
	}()

	var token *jwt.Token
	token, err = s.parseToken(ctx, tokenString)
	if err != nil {
		return tokenStr, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if expiresAt, ok := claims["exp"].(float64); ok &&
			time.Unix(int64(expiresAt), 0).Sub(time.Now()) > time.Duration(configs.Conf.Timeouts.DefaultTimeout)*time.Second {
			return tokenStr, fmt.Errorf("token expites at time is more then default timeout value")
		}
		return s.GenerateToken(ctx)
	}

	return tokenStr, jwt.ErrSignatureInvalid
}
