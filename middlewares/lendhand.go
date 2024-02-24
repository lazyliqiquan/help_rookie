package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
)

// 需要登录
func PublishLendHand() gin.HandlerFunc {
	return func(c *gin.Context) {
		publishLendHandBan, err := models.RDB.Get(c, "publishLendHandBan").Result()
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
		if publishLendHandBan != "permit" && !models.UserPermit(models.Root, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		// 局部判断
		if !models.UserPermit(models.PublishLendHand, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to lend hand",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func EditLendHand() gin.HandlerFunc {
	return func(c *gin.Context) {
		editLendHandBan, err := models.RDB.Get(c, "editLendHandBan").Result()
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
		if editLendHandBan != "permit" && !models.UserPermit(models.Root, userBan) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		// 局部判断
		if !models.UserPermit(models.EditLendHand, userBan) {
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
func ViewLendHand() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewLendHandBan, err := models.RDB.Get(c, "viewLendHandBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
			return
		}
		loginViewLendHandBan, err := models.RDB.Get(c, "loginViewLendHandBan").Result()
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
		if viewLendHandBan != "permit" && (userId == 0 || !models.UserPermit(models.Root, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginViewLendHandBan != "permit" && (userId == 0 || !models.UserPermit(models.ViewLendHand, userBan)) {
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
