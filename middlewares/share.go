package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
)

func ViewShareCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewShareCodeBan, err := models.RDB.Get(c, "viewShareCodeBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
			return
		}
		loginViewShareCodeBan, err := models.RDB.Get(c, "loginViewShareCodeBan").Result()
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
		if viewShareCodeBan != "permit" && (userId == 0 || !models.UserPermit(models.Root, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginViewShareCodeBan != "permit" && (userId == 0 || !models.UserPermit(models.ViewShareCode, userBan)) {
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
func PublishShareCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		PublishShareCodeBan, err := models.RDB.Get(c, "publishShareCodeBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
			return
		}
		loginPublishShareCodeBan, err := models.RDB.Get(c, "loginPublishShareCodeBan").Result()
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
		if PublishShareCodeBan != "permit" && (userId == 0 || !models.UserPermit(models.Root, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginPublishShareCodeBan != "permit" && (userId == 0 || !models.UserPermit(models.PublishShareCode, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to post a share code",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
