package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/models"
	"gorm.io/gorm"
)

// @Tags 公共方法
// @Summary 请求求助列表
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Param baseOffset formData string true "baseOffset"
// @Param size formData string true "size"
// @Param sortOption formData string true "sortOption"
// @Param language formData string true "language"
// @Param status formData string true "status"
// @Success 200 {string} json "{"code":"0"}"
// @Router /seek-help-list [post]
func RequestSeekHelpList(c *gin.Context) {
	baseOffsetStr := c.PostForm("baseOffset")
	sizeStr := c.PostForm("size")
	sortOption := c.PostForm("sortOption")
	language := c.PostForm("language")
	status := c.PostForm("status")
	if helper.IsNuiStrs(baseOffsetStr, sizeStr, sortOption, language, status) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	baseOffset, err := strconv.Atoi(baseOffsetStr)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "BaseOffset not an integer",
		})
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Size not an integer",
		})
		return
	}
	var total int64
	tx := models.DB.Model(&models.SeekHelps{})
	if language != "All" {
		tx = tx.Where("language = ?", language)
	}
	if status == "1" {
		tx = tx.Where("status = ?", true)
	} else if status == "2" {
		tx = tx.Where("status = ?", false)
	}
	if err := tx.Count(&total).Error; err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	// 初始化，防止前端收到null,而不是[]
	seekHelps := make([]models.SeekHelps, 0)
	if sortOption == "0" {
		tx = tx.Order("id DESC")
	} else {
		tx = tx.Order("reward DESC")
	}
	// todo 算了，把后面三个筛选条件也给实现了吧，preload结合手写sort
	if err := tx.Offset(baseOffset).Limit(size).Find(&seekHelps).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failure",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Mysql operation failure",
	})
}