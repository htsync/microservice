package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
		"app.log",
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
