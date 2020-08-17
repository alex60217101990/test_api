package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alex60217101990/test_api/internal/banner"
	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/encrypt"
	"github.com/alex60217101990/test_api/internal/encrypt/jwt"
	fast_http "github.com/alex60217101990/test_api/internal/http-server/fast-http"
	"github.com/alex60217101990/test_api/internal/logger"
	"github.com/alex60217101990/test_api/internal/repository"
	"github.com/alex60217101990/test_api/internal/repository/postgres"
)

var (
	confFile     = flag.String("conf", "../../deploy/configs/application.yaml", "Config file path")
	isDebug      = flag.Bool("debug", false, "Is Debug")
	genKeys      = flag.Bool("gk", false, "Generate new auth RSA keys pair ?")
	genKeysForce = flag.Bool("f", false, "Generate new auth RSA keys pair with force option ?")
	loggerType   configs.LoggerType
)

func main() {
	flag.Var(&loggerType, "lt", "Type of logger usage")
	flag.Parse()

	// Load configs file
	err := configs.ReadConfigFile(*confFile)
	if err != nil {
		_, _ = logger.CmdError.Println(err)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	// Init loggers
	logger.InitLoggerSettings()

	configs.Conf.IsDebug = *isDebug

	if configs.Conf.IsDebug {
		// Showing startup banner
		banner.StartUpBanner()
	}

	if *genKeys {
		currentDir, err := os.Getwd()
		if err != nil {
			logger.AppLogger.Fatal(err)
		}
		err = encrypt.InitKeys("test-api.rsa", currentDir, *genKeysForce)
		if err != nil {
			logger.AppLogger.Error(err)
		}
	}

	// Run repository connections
	var repo repository.Repository
	repoCtx, repoCtxCloser := context.WithCancel(context.Background())

	switch configs.Conf.DB.RepoType {
	case configs.RepoPostgres:
		repo = postgres.NewPostgresRepository()
	default:
		logger.AppLogger.Fatal(fmt.Errorf("config error: invalid repository type %v", configs.Conf.DB.RepoType))
	}
	err = repo.Connect(repoCtx)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	// Init secret service
	secret := jwt.NewSecretService(context.Background(), repo)

	// Run Http/s server
	server := fast_http.NewFastHttpServer(repo, secret)
	server.Init()
	go server.Run()

	logger.CmdServer.Printf("🚀 %s service started...\n", strings.ToUpper(configs.Conf.ServiceName))

	var Stop = make(chan os.Signal, 1)
	signal.Notify(Stop,
		syscall.SIGTERM,
		syscall.SIGINT,
		// syscall.SIGKILL,
		syscall.SIGABRT,
	)
	for range Stop {
		server.Close()
		repo.Close()
		repoCtxCloser()
		logger.CmdServer.Printf("🚫 %s service stoped...\n", strings.ToUpper(configs.Conf.ServiceName))
		return
	}
}
