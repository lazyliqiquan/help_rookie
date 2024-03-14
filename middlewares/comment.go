package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/models"
)

func ViewComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		viewCommentBan, err := models.RDB.Get(c, "viewCommentBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
		}
		loginViewCommentBan, err := models.RDB.Get(c, "loginViewCommentBan").Result()
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
		if viewCommentBan != "permit" && (userId == 0 || !models.JudgePermit(models.Admin, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
		}
		if loginViewCommentBan != "permit" && (userId == 0 || !models.JudgePermit(models.Login, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You are not logged in or do not have browsing rights",
			})
			c.Abort()
		}
		c.Next()
	}
}
func PublishComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		PublishCommentBan, err := models.RDB.Get(c, "publishCommentBan").Result()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			c.Abort()
			return
		}
		loginPublishCommentBan, err := models.RDB.Get(c, "loginPublishCommentBan").Result()
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
		if PublishCommentBan != "permit" && (userId == 0 || !models.JudgePermit(models.Login, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The website is currently in safe mode and can only be operated by administrators",
			})
			c.Abort()
			return
		}
		if loginPublishCommentBan != "permit" && (userId == 0 || !models.JudgePermit(models.PublishComment, userBan)) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You do not have permission to post a comment",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
