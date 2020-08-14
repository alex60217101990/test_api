package main

import (
	"flag"
	"os"
	"syscall"

	"github.com/alex60217101990/test_api/internal/banner"
	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/encrypt"
	"github.com/alex60217101990/test_api/internal/logger"
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
	// Showing startup banner
	banner.StartUpBanner()

	// Load configs file
	err := configs.ReadConfigFile(*confFile)
	if err != nil {
		logger.CmdError.Println(err)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}

	// Init loggers
	logger.InitLoggerSettings()

	if *genKeys {
		currentDir, err := os.Getwd()
		if err != nil {
			logger.AppLogger.Fatal(err)
		}
		encrypt.InitKeys("test-api.rsa", currentDir, *genKeysForce)
	}
}
