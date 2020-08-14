package fast_http

import (
	server "github.com/alex60217101990/test_api/internal/http-server"
	"github.com/valyala/fasthttp"
)

type FastHttpServer struct {
	//...
}

func NewFastHttpServer() server.Server {
	return new(FastHttpServer)
}

func (s *FastHttpServer) addRoute(routeType server.RouteType, path string, handler fasthttp.RequestHandler) {

}
