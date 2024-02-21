package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/models"
)

// @Tags 用户方法
// @Summary 新添加求助
// @Accept multipart/form-data
// @Param file formData file true "File to upload"
// @Param document formData string true "document"
// @Param language formData string true "language"
// @Param score formData string true "score"
// @Param createTime formData string true "createTime"
// @Success 200 {string} json "{"code":"0"}"
// @Router /add-seek-help [post]
func AddSeekHelp(c *gin.Context) {
	document := c.PostForm("document")
	language := c.PostForm("language")
	createTime := c.PostForm("createTime")
	scoreStr := c.PostForm("score")
	if helper.IsNuiStrs(document, language, createTime, scoreStr) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The uploaded file cannot be obtained",
		})
		return
	}
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "score cannot be parsed into a number",
		})
		return
	}
	filePath := config.Config.CodeFilePath + helper.GetUUID() + ".txt"
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "File saving failure",
		})
		return
	}
	seekHelp := &models.SeekHelps{
		Document:   document,
		Score:      score,
		CreateTime: createTime,
		CodePath:   filePath,
		Language:   language,
	}
	// models.DB.Model(&models.SeekHelp{}).Create()
}
