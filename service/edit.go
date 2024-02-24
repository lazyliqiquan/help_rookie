package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/models"
	"gorm.io/gorm"
)

// @Tags 编辑方法
// @Summary 添加求助
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
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
	logger.Infoln(filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "File saving failure",
		})
		return
	}
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Unable to get user id",
		})
		return
	}
	seekHelp := &models.SeekHelps{
		Document:   document,
		Score:      score,
		CreateTime: createTime,
		CodePath:   filePath,
		Language:   language,
		UsersID:    userId,
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.SeekHelps{}).Create(seekHelp).Error
		if err != nil {
			return err
		}
		var userScore int
		err = tx.Model(&models.Users{}).
			Select("score").Where("id = ?", userId).Scan(&userScore).Error
		if err != nil {
			return err
		}
		if userScore < score || score <= 0 {
			return fmt.Errorf("score error")
		}
		userScore -= score
		return tx.Model(&models.Users{}).
			Where("id = ?", userId).Update("score", userScore).Error
	})
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Add seek help",
	})
}

// @Tags 编辑
// @Summary 修改求助
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Param file formData file true "File to upload"
// @Param document formData string true "document"
// @Param language formData string true "language"
// @Param score formData string true "score"
// @Param createTime formData string true "createTime"
// @Success 200 {string} json "{"code":"0"}"
// @Router /add-seek-help [post]
func EditSeekHelp(c *gin.Context) {
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
	logger.Infoln(filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "File saving failure",
		})
		return
	}
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Unable to get user id",
		})
		return
	}
	seekHelp := &models.SeekHelps{
		Document:   document,
		Score:      score,
		CreateTime: createTime,
		CodePath:   filePath,
		Language:   language,
		UsersID:    userId,
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.SeekHelps{}).Create(seekHelp).Error
		if err != nil {
			return err
		}
		var userScore int
		err = tx.Model(&models.Users{}).
			Select("score").Where("id = ?", userId).Scan(&userScore).Error
		if err != nil {
			return err
		}
		if userScore < score || score <= 0 {
			return fmt.Errorf("score error")
		}
		userScore -= score
		return tx.Model(&models.Users{}).
			Where("id = ?", userId).Update("score", userScore).Error
	})
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Add seek help",
	})
}

// @Tags 编辑
// @Summary 浏览求助
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Param file formData file true "File to upload"
// @Param document formData string true "document"
// @Param language formData string true "language"
// @Param score formData string true "score"
// @Param createTime formData string true "createTime"
// @Success 200 {string} json "{"code":"0"}"
// @Router /add-seek-help [post]
func ViewSeekHelp(c *gin.Context) {
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
	logger.Infoln(filePath)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "File saving failure",
		})
		return
	}
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Unable to get user id",
		})
		return
	}
	seekHelp := &models.SeekHelps{
		Document:   document,
		Score:      score,
		CreateTime: createTime,
		CodePath:   filePath,
		Language:   language,
		UsersID:    userId,
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.SeekHelps{}).Create(seekHelp).Error
		if err != nil {
			return err
		}
		var userScore int
		err = tx.Model(&models.Users{}).
			Select("score").Where("id = ?", userId).Scan(&userScore).Error
		if err != nil {
			return err
		}
		if userScore < score || score <= 0 {
			return fmt.Errorf("score error")
		}
		userScore -= score
		return tx.Model(&models.Users{}).
			Where("id = ?", userId).Update("score", userScore).Error
	})
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Add seek help",
	})
}
