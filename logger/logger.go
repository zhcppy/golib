package logger

import (
	"github.com/zhcppy/golib/common/util"
	"go.uber.org/zap"
)

// default logger
var L *zap.SugaredLogger

func init() {
	L = New()
}

func New(name ...string) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	if util.RunEnv() == util.EnvProd {
		config = zap.NewProductionConfig()
	}
	logger, err := config.Build()
	if err != nil {
		panic(err.Error())
	}
	if len(name) > 0 {
		logger.Named(name[0])
	}
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}
