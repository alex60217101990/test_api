package fast_http

import (
	"encoding/json"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) Signin(ctx *fasthttp.RequestCtx) {
	var creds models.Credentials
	// Get the JSON body and decode into credentials
	err := json.Unmarshal(ctx.PostBody(), &creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	
}
