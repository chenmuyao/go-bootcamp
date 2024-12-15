package ioc

import (
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() logger.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
