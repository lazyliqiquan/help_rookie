package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	_ "github.com/lazyliqiquan/help_rookie/docs"
	"github.com/lazyliqiquan/help_rookie/middlewares"
	"github.com/lazyliqiquan/help_rookie/service"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(config *config.WebConfig) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 公共方法
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	// 用户方法
	// 管理员方法
	// 超级管理员方法
	httpsPath := config.WebPath
	if config.Debug {
		httpsPath = config.DebugWebPath
	}
	err := r.RunTLS(httpsPath, "assets/https/cert.pem", "assets/https/key.pem")
	if err != nil {
		log.Fatal(err)
	}
	return r
}
