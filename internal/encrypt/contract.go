package encrypt

import "context"

type SecretService interface {
	Init(ctx context.Context) (err error)
	GenerateToken(ctx context.Context) (token string, err error)
	VerifyTokenString(ctx context.Context, tokenString string) (newCtx context.Context, err error)
	RefreshToken(ctx context.Context, tokenString string) (tokenStr string, err error)
}
