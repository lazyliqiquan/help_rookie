package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
)

// 需要登录
func PublishSeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		publishSeekHelpBan, err := models.RDB.Get(c, "publishSeekHelpBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		userBan := c.GetInt("ban")
		// 全局判断
		if publishSeekHelpBan != "permit" && !models.JudgePermit(models.Admin, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		// 局部判断
		if !models.JudgePermit(models.PublishSeekHelp, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to seek help",
			})
			c.Abort()
		}
		logger.Infoln("PublishSeekHelp")
		c.Next()
	}
}

func EditSeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		editSeekHelpBan, err := models.RDB.Get(c, "editSeekHelpBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		userBan := c.GetInt("ban")
		// 全局判断
		if editSeekHelpBan != "permit" && !models.JudgePermit(models.Admin, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		// 局部判断
		if !models.JudgePermit(models.EditSeekHelp, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to modify",
			})
			c.Abort()
		}
		logger.Infoln("EditSeekHelp")
		c.Next()
	}
}

// 可能需要登录
func ViewSeekHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewSeekHelpBan, err := models.RDB.Get(c, "viewSeekHelpBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		loginViewSeekHelpBan, err := models.RDB.Get(c, "loginViewSeekHelpBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		// 全局判断
		if viewSeekHelpBan != "permit" && (userId == 0 || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		if loginViewSeekHelpBan != "permit" && (userId == 0 || !models.JudgePermit(models.Login, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not logged in or do not have browsing rights",
			})
			c.Abort()
		}
		logger.Infoln("ViewSeekHelp")
		c.Next()
	}
}
