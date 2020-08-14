package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(jsonEncoder, consoleDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	return &ZapLogger{
		logger: zap.New(
			core,
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	}
}

func (l *ZapLogger) GetNativeLogger() interface{} {
	return l.logger
}

func (l *ZapLogger) Close() {
	defer l.logger.Sync()
}

func (l *ZapLogger) Info(msg string) {
	defer func() {
		_ = l.logger.Sync()
	}()

	l.logger.Info(msg)
}

func (l *ZapLogger) Infof(tpl string, args map[string]interface{}) {
	defer func() {
		_ = l.logger.Sync()
	}()

	fields := make([]zapcore.Field, 0)
	for k, v := range args {
		fields = append(fields, zap.Any(k, v))
	}
	l.logger.Info(tpl, fields...)
}

func (l *ZapLogger) Error(err error) {
	defer func() {
		_ = l.logger.Sync()
	}()

	l.logger.Error(err.Error())
}

func (l *ZapLogger) Errorf(tpl string, args map[string]interface{}) {
	defer func() {
		_ = l.logger.Sync()
	}()

	fields := make([]zapcore.Field, 0)
	for k, v := range args {
		fields = append(fields, zap.Any(k, v))
	}
	l.logger.Error(tpl, fields...)
}

func (l *ZapLogger) Warn(msg string) {
	defer func() {
		_ = l.logger.Sync()
	}()

	l.logger.Warn(msg)
}

func (l *ZapLogger) Warnf(tpl string, args map[string]interface{}) {
	defer func() {
		_ = l.logger.Sync()
	}()

	fields := make([]zapcore.Field, 0)
	for k, v := range args {
		fields = append(fields, zap.Any(k, v))
	}
	l.logger.Warn(tpl, fields...)
}

func (l *ZapLogger) Fatal(err error) {
	defer func() {
		_ = l.logger.Sync()
	}()

	l.logger.Fatal(err.Error())
}

func (l *ZapLogger) Fatalf(tpl string, args map[string]interface{}) {
	defer func() {
		_ = l.logger.Sync()
	}()

	fields := make([]zapcore.Field, 0)
	for k, v := range args {
		fields = append(fields, zap.Any(k, v))
	}

	l.logger.Fatal(tpl, fields...)
}
