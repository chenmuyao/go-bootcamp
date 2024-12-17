package startup

import (
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	l, _ := zap.NewDevelopment()
	return logger.NewZapLogger(l)
}
