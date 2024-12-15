package startup

import "github.com/chenmuyao/go-bootcamp/pkg/logger"

func InitLogger() logger.Logger {
	return logger.NewNopLogger()
}
