package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Init(loggerInstance *zap.SugaredLogger) {
	logger = loggerInstance
}

// 获取网站配置信息
func GetWebConfig(c *gin.Context) {
	models.RDB.Get(c, "sd").Int()
}
