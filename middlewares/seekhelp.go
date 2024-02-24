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
			return
		}
		userBan := c.GetInt("ban")
		// 全局判断
		if publishSeekHelpBan != "permit" && !models.UserPermit(models.Root, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		// 局部判断
		if !models.UserPermit(models.PublishSeekHelp, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to seek help",
			})
			c.Abort()
			return
		}
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
			return
		}
		userBan := c.GetInt("ban")
		// 全局判断
		if editSeekHelpBan != "permit" && !models.UserPermit(models.Root, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		// 局部判断
		if !models.UserPermit(models.PublishSeekHelp, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to modify",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 也许不需要登录
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
			return
		}
		loginViewSeekHelpBan, err := models.RDB.Get(c, "loginViewSeekHelpBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
			return
		}
		userId := c.GetInt("id")
		userBan := c.GetInt("ban")
		// 全局判断
		if viewSeekHelpBan != "permit" && (userId == 0 || !models.UserPermit(models.Root, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginViewSeekHelpBan != "permit" && (userId == 0 || !models.UserPermit(models.ViewSeekHelp, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not logged in or do not have browsing rights",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
