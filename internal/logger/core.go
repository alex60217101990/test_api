package logger

import (
	"log"
	"os"
	"sync"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/fatih/color"
)

var (
	PackageOnceLoad sync.Once

	CmdServer = color.New(color.FgHiGreen, color.Bold)
	CmdError  = color.New(color.FgHiRed, color.Bold)
	CmdInfo   = color.New(color.FgHiBlue, color.Faint)

	AppLogger  Logger
	RepoLogger *log.Logger
)

func InitLoggerSettings() {
	PackageOnceLoad.Do(func() {
		CmdServer.Println("Run once - 'logger' package loading...")

		RepoLogger = log.New(os.Stderr, "[REPO] ", log.LstdFlags|log.Lshortfile|log.LUTC)

		switch configs.Conf.LoggerType {
		case configs.Base:
			// TODO: standart golang logger implementation...
		case configs.Zap:
			AppLogger = NewZapLogger()
		default:
			log.Fatalf("invalid logger type: %+v\n", configs.Conf.LoggerType)
		}
	})
}

func CloseLoggers() {
	AppLogger.Close()

	CmdServer.Println("'logger' package stoped...")
}
