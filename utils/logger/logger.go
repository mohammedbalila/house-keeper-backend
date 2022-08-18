package logger

import (
	"strings"

	"github.com/mustafabalila/golang-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger JSON based STDOUT logger.
var logger *zap.Logger

// GetLoggerInstance returns, or creates an instance of the logger
func GetLoggerInstance() *zap.Logger {
	if logger == nil {
		cfg := config.GetConfig()
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.MessageKey = "message"
		config.DisableStacktrace = true

		switch strings.ToLower(cfg.LogLevel) {
		case "silent":
		case "trace":
		case "debug":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
				break
			}
		case "info":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
				break
			}
		case "warn":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
				break
			}
		case "error":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
				break
			}
		case "fatal":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
				break
			}
		case "default":
			{
				config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
				break
			}
		}
		config.InitialFields = map[string]interface{}{
			"service": cfg.ServiceName,
		}

		log, err := config.Build()
		if err != nil {
			panic(err)
		}

		logger = log
	}

	return logger
}
