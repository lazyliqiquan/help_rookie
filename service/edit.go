package service

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/models"
	"gorm.io/gorm"
)

var (
	editLimitKeys = []string{"maxDocumentHeight", "maxDocumentLength", "maxPictureSize", "maxCodeFileSize"}
)

// 这个函数就不用swagger了,文件以及document都是需要的，并且有相关性
func AddSeekHelp(c *gin.Context) {
	_imageNum := c.PostForm("imageNum")
	_reward := c.PostForm("reward")
	document := c.PostForm("document")
	language := c.PostForm("language")
	uploadTime := c.PostForm("uploadTime")
	tags := c.PostForm("tags")
	if helper.IsNuiStrs(_imageNum, _reward, document, language, uploadTime) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Parameter null",
		})
		return
	}
	imageNum, err := strconv.Atoi(_imageNum)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The imageNum parameter is not an integer",
		})
		return
	}
	reward, err := strconv.Atoi(_reward)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The reward parameter is not an integer",
		})
		return
	}
	if reward <= 0 || reward > 9 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The reward parameter invalid",
		})
		return
	}
	userId := c.GetInt("id")
	user := models.Users{}
	err = models.DB.Model(&models.Users{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	if reward > user.Reward {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "User reward amount is insufficient",
		})
		return
	}
	imageFiles := []multipart.File{}
	for i := 0; i < imageNum; i++ {
		file, _, err := c.Request.FormFile("image" + strconv.Itoa(i))
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Image file parsing error",
			})
			return
		}
		// 这里关闭的文件应该不总是最后一个吧(如果程序内存溢出可以考虑文件是否及时关闭)
		defer file.Close()
		imageFiles = append(imageFiles, file)
	}
	codeFile, _, err := c.Request.FormFile("codeFile")
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Code file parsing error",
		})
		return
	}
	codeFilePath := config.Config.CodeFilePath + helper.GetUUID() + ".txt"
	imageFilesPath := []string{}
	for range imageFiles {
		// 反正存的是二进制，具体的文件类型问题应该不大吧
		imageFilesPath = append(imageFilesPath, config.Config.CodeFilePath+helper.GetUUID())
	}
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 不懂没有seekHelpId直接追加关联可不可以
		err = tx.Model(&models.Users{ID: userId}).Association("SeekHelps").Append(
			&models.SeekHelps{
				Reward:     reward,
				CreateTime: uploadTime,
				CodePath:   codeFilePath,
				Language:   language,
				Document:   document,
				ImagePath:  imageFilesPath,
				Tags:       tags,
			},
		)
		if err != nil {
			return err
		}
		err = tx.Model(&models.Users{ID: userId}).Update("reward", user.Reward-reward).Error
		if err != nil {
			return err
		}
		for i, v := range imageFilesPath {
			err = helper.SaveAFile(v, imageFiles[i])
			if err != nil {
				return err
			}
		}
		return helper.SaveAFile(codeFilePath, codeFile)
	})
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql transaction error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Add seek help successfully",
	})
}

// @Tags 用户方法
// @Summary 检测编辑权限并获取环境配置
// @Accept multipart/form-data
//
//	@Param Authorization header string true	"Authentication header"
//
// @Param editOption formData string true "editOption"
// @Param seekHelpId formData string false "seekHelpId"
// @Param lendHandId formData string false "lendHandId"
// @Success 200 {string} json "{"code":"0"}"
// @Router /preEdit [post]
func PreEdit(c *gin.Context) {
	// 这里一般都会获取得到
	editOption := c.GetInt("editOption")
	userId := c.GetInt("id")
	var remainReward int
	var language string
	seekHelp := &models.SeekHelps{}
	lendHand := &models.LendHands{}
	var err error
	if editOption == 0 { //新增求助
		// 用户的金额是否大于0
		if err := models.DB.Model(&models.Users{ID: userId}).Select("reward").Scan(&remainReward).Error; err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			return
		}
		if remainReward <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "You have no amount to issue a reward",
			})
			return
		}
	}
	// 需要用到seekHelpId的情况
	if editOption == 1 || editOption == 2 {
		seekHelp.ID, err = strconv.Atoi(c.PostForm("seekHelpId"))
		if err != nil {
			logger.Errorln(err)
			// 不是整数的情况应该交给前端处理，我们不需要额外说明
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Seek help id not an integer",
			})
			return
		}
		err = models.DB.Model(&models.SeekHelps{ID: seekHelp.ID}).Select("status", "ban", "lend_hand_sum", "language").First(seekHelp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { //传递过来的seekHelpId不存在，用户输入的url有问题
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Url error : seek help id not exist",
				})
				return
			}
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			return
		}
		language = seekHelp.Language
	}
	// 需要用到lendHandId的情况
	if editOption == 3 {
		lendHand.ID, err = strconv.Atoi(c.PostForm("lendHandId"))
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Lend hand id not an integer",
			})
			return
		}
		err = models.DB.Model(&models.LendHands{ID: lendHand.ID}).Preload("SeekHelps").Select("status", "ban").First(lendHand).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { //传递过来的seekHelpId不存在，用户输入的url有问题
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Url error : lend hand id not exist",
				})
				return
			}
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			return
		}
		language = lendHand.SeekHelps.Language
	}
	if editOption == 1 { //编辑求助
		// 是否已经有帮助帖子
		// 是否被管理员禁止修改
		if seekHelp.LendHandSum > 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "A user has posted a help post that you cannot modify",
			})
			return
		}
		if !models.JudgePermit(models.SEdit, seekHelp.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Your post has been disabled by the administrator",
			})
			return
		}
	} else if editOption == 2 { //新增帮助
		// 用户是否已经发布过针对该求助帖子的帮助帖子
		// 该求助帖子是否已经被解决
		// 求助帖子是否被管理员禁止添加帮助帖子
		tempUser := &models.Users{}
		// todo 有没有可能默认情况下seekHelpId也是需要preload SeekHelps才会查询,这样会导致查询失败
		err = models.DB.Model(&models.Users{ID: userId}).Preload("LendHands", "seek_help_id = ?", seekHelp.ID).Select("id").First(tempUser).Error
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failed",
			})
			return
		}
		// todo 不一定能正常工作，Preload不太熟
		if len(tempUser.LendHands) > 0 { //说明用户已经对该求助帖子提交过帮助了
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Only one help post is allowed per help post",
			})
			return
		}
		if seekHelp.Status {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The issue has been resolved and you are unable to post help",
			})
			return
		}
		if !models.JudgePermit(models.SAddLendHand, seekHelp.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The post has been disabled by the administrator",
			})
			return
		}
	} else if editOption == 3 { //修改帮助
		// 是否已经被求助者接受
		// 是否被管理员禁止修改
		if lendHand.Status {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Your solution has been accepted by the supplicant and cannot be modified",
			})
			return
		}
		if !models.JudgePermit(models.LEdit, lendHand.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The post has been disabled by the administrator",
			})
			return
		}
	}
	editLimit := []int{}
	for _, v := range editLimitKeys {
		limit, err := models.RDB.Get(c, v).Int()
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failed",
			})
			return
		}
		editLimit = append(editLimit, limit)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "The edit configuration was obtained successfully",
		"data": gin.H{
			"documentLimit": editLimit,
			"remainReward":  remainReward,
			"language":      language,
		},
	})
}
