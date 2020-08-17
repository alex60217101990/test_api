package fast_http

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) Ping(ctx *fasthttp.RequestCtx) {
	err := s.repo.Ping(ctx)
	if err != nil {
		ctx.Error(http.StatusText(fasthttp.StatusServiceUnavailable), fasthttp.StatusServiceUnavailable)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
