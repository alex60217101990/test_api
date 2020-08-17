package fast_http

import (
	"fmt"
	"time"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"

	"github.com/alex60217101990/test_api/internal/encrypt"
	server "github.com/alex60217101990/test_api/internal/http-server"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/fasthttp/router"

	"github.com/valyala/fasthttp"
)

type FastHttpServer struct {
	repo   repository.Repository
	sicret encrypt.SecretService
	server *fasthttp.Server
}

func NewFastHttpServer(repo repository.Repository, sicret encrypt.SecretService) server.Server {
	return &FastHttpServer{
		repo:   repo,
		sicret: sicret,
	}
}

func (s *FastHttpServer) routing() *router.Router {
	r := router.New()

	r.GET("/liveness", s.Ping)
	r.GET("/readiness", s.Ping)

	grAuth := r.Group("/v1/auth")
	grAuth.POST("/singup", s.Singup)
	grAuth.POST("/singin", s.Signin)
	grAuth.POST("/refresh", s.RefreshToken)

	grProduct := r.Group("/v1/products")
	grProduct.GET("/single", s.AuthMiddleware(s.SingleProduct))
	grProduct.GET("/list", s.AuthMiddleware(s.ProductsList))
	grProduct.PUT("/create", s.AuthMiddleware(s.CreateProduct))
	grProduct.DELETE("/delete", s.AuthMiddleware(s.DeleteProduct))
	grProduct.POST("/update", s.AuthMiddleware(s.UpdateProduct))

	grCategory := r.Group("/v1/categories")
	grCategory.GET("/single", s.AuthMiddleware(s.SingleCategory))
	grCategory.GET("/list", s.AuthMiddleware(s.CategoriesList))
	grCategory.PUT("/create", s.AuthMiddleware(s.CreateCategory))
	grCategory.DELETE("/delete", s.AuthMiddleware(s.DeleteCategory))
	grCategory.POST("/update", s.AuthMiddleware(s.UpdateCategory))

	grRelations := r.Group("/v1/relations")
	grRelations.POST("/add-prod-cat-relation", s.AuthMiddleware(s.AddCategoryProductRelation))
	grRelations.DELETE("/del-prod-cat-relation", s.AuthMiddleware(s.DeleteCategoryProductRelation))

	return r
}

func (s *FastHttpServer) Init() {
	s.server = &fasthttp.Server{
		Name:               configs.Conf.ServiceName,
		Concurrency:        100000,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 3 * time.Second,
		ReadBufferSize:     1 << 10,
		WriteBufferSize:    1 << 10,
		ReadTimeout:        7 * time.Second,
		WriteTimeout:       7 * time.Second,
		IdleTimeout:        15 * time.Second,
		Logger:             logger.AppLogger,
		Handler:            fasthttp.CompressHandler(s.PanicMiddleware(s.CorsMiddleware(s.routing().Handler))),
	}

	if configs.Conf.IsDebug {
		s.server.LogAllErrors = true
	}
}

func (s *FastHttpServer) Run() {
	// You can check the access using openssl command:
	// $ openssl s_client -connect localhost:8080 << EOF
	// > GET /
	// > Host: localhost
	// > EOF
	//
	// $ openssl s_client -connect localhost:8080 << EOF
	// > GET /
	// > Host: 127.0.0.1:8080
	// > EOF
	//

	// preparing first host
	// cert, priv, err := encrypt.GenerateCert(fmt.Sprintf("%s:%d", configs.Conf.Server.Host, configs.Conf.Server.Port))
	// if err != nil {
	// 	logger.AppLogger.Fatal(err)
	// }
	// err = s.server.AppendCertEmbed(cert, priv)
	// if err != nil {
	// 	logger.AppLogger.Fatal(err)
	// }
	// cert, priv, err = encrypt.GenerateCert("127.0.0.1:8080")
	// if err != nil {
	// 	logger.AppLogger.Fatal(err)
	// }
	// err = s.server.AppendCertEmbed(cert, priv)
	// if err != nil {
	// 	logger.AppLogger.Fatal(err)
	// }

	logger.CmdInfo.Println("🔥 FataHttp server started.")
	logger.CmdError.Println(s.server.ListenAndServe(
		fmt.Sprintf("%s:%d", configs.Conf.Server.Host, configs.Conf.Server.Port)))
	// fmt.Println(s.server.ListenAndServeTLS(":8080",
	// 	"/Users/aleksandr/generate-ssl-certs-for-local-development/your-certs/local.dev.crt",
	// 	"/Users/aleksandr/generate-ssl-certs-for-local-development/your-certs/local.dev.key"))
}

func (s *FastHttpServer) Close() error {
	defer logger.CmdInfo.Println("🔥 FataHttp server stoped.")
	return s.server.Shutdown()
}
