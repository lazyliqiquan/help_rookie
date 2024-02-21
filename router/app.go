package router

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
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
	if config.Debug {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Access-Control-Allow-Origin"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	r.MaxMultipartMemory = config.MaxMultipartMemory << 20
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 公共方法
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	r.POST("/find-password", service.FindPassword)
	// 用户方法
	userAuth := r.Group("/", middlewares.AuthCheck(middlewares.UserLevel))
	userAuth.POST("/add-seek-help", service.AddSeekHelp)
	// 管理员方法
	adminAuth := r.Group("/", middlewares.AuthCheck(middlewares.UserLevel))
	adminAuth.POST("/delete-user", service.DeleteUser)
	// 超级管理员方法
	// rootAuth := r.Group("/root", middlewares.AuthCheck(middlewares.UserLevel))

	if !config.Debug {
		err := r.RunTLS(config.WebPath, "assets/https/cert.pem", "assets/https/key.pem")
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}
