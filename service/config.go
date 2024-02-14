package service

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func ServiceInit(loggerInstance *zap.SugaredLogger) {
	logger = loggerInstance
}
