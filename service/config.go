package service

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func Init(loggerInstance *zap.SugaredLogger) {
	logger = loggerInstance
}
