package fast_http

import (
	"strings"

	"github.com/alex60217101990/test_api/internal/helpers"
	"github.com/alex60217101990/test_api/internal/models"
	"github.com/alex60217101990/test_api/internal/repository"

	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) Singup(ctx *fasthttp.RequestCtx) {
	creeds, err := parseCreeds(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}
	err = s.repo.InsertUser(ctx, new(models.User).FromCreeds(creeds))
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *FastHttpServer) RefreshToken(ctx *fasthttp.RequestCtx) {
	var err error
	tokenStr := strings.TrimPrefix(string(ctx.Request.Header.Peek("Authorization")), "Bearer ")
	tokenStr, err = s.sicret.RefreshToken(ctx, tokenStr)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
	messagePrint(ctx, struct {
		NewToken string `json:"new_token"`
	}{
		NewToken: tokenStr,
	}, fasthttp.StatusOK)
}

func (s *FastHttpServer) Signin(ctx *fasthttp.RequestCtx) {
	creeds, err := parseCreeds(ctx)
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	var user *models.User
	user, err = s.repo.GetUserByCreeds(ctx, creeds, struct{}{})
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusUnauthorized)
		return
	}
	user.ID = 0

	var token string
	token, err = s.sicret.GenerateToken(helpers.AddToContext(ctx, repository.UserSessionKey, user))
	if err != nil {
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	messagePrint(ctx, struct {
		User  *models.User
		Token string
	}{
		User:  user,
		Token: token,
	}, fasthttp.StatusOK)
}
