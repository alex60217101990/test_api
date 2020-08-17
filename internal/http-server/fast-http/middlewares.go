package fast_http

import (
	"fmt"
	"strings"

	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func (s *FastHttpServer) CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func (s *FastHttpServer) PanicMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					logger.AppLogger.Error(fmt.Errorf("panic middleware detect, %w", err))
					errorPrint(ctx, err, fasthttp.StatusInternalServerError)
				}
			}
			ctx.Done()
		}()
		next(ctx)
	}
}

func (s *FastHttpServer) AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenStr := strings.TrimPrefix(string(ctx.Request.Header.Peek("Authorization")), "Bearer ")
		if len(tokenStr) == 0 {
			errorPrint(ctx, fmt.Errorf("token value is empty"), fasthttp.StatusBadRequest)
			return
		}

		_, err := s.sicret.VerifyTokenString(ctx, tokenStr)
		if err != nil {
			errorPrint(ctx, err, fasthttp.StatusBadRequest)
			return
		}

		next(ctx)
	}
}
