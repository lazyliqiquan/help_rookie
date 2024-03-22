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
	r.MaxMultipartMemory = int64(config.MaxMultipartMemory) << 20
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 登录
	r.POST("/login", service.Login)
	// 启动安全模式后，谁都不可以使用这些功能
	anyoneAuth := r.Group("/", middlewares.OtherSafeModel())
	anyoneAuth.POST("/send-code", service.SendCode)
	anyoneAuth.POST("/register", service.Register)
	anyoneAuth.POST("/find-password", service.FindPassword)
	// 启动安全模式后，仅管理员可用
	managerAuth := r.Group("/", middlewares.TokenSafeModel())
	managerAuth.POST("/seek-help-list", service.RequestSeekHelpList)

	// todo 需要登录的操作(应该不会影响到之前的方法吧)
	loginAuth := managerAuth.Group("/", middlewares.LoginModel())
	// 需要编辑权限的操作
	editAuth := loginAuth.Group("/", middlewares.JudgeEdit())
	editAuth.POST("/preEdit", service.PreEdit)
	editAuth.POST("/add-seek-help", service.AddSeekHelp)

	if !config.Debug {
		err := r.RunTLS(config.WebPath, "assets/https/cert.pem", "assets/https/key.pem")
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}
