package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
)

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
		if viewLendHandBan != "permit" && (userId == 0 || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginViewLendHandBan != "permit" && (userId == 0 || !models.JudgePermit(models.Login, userBan)) {
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
