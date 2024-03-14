package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/middlewares"
	"github.com/lazyliqiquan/help_rookie/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SendCode
// @Tags 公共方法
// @Summary 发送验证码(一个验证码只能处理一个操作，用完就要删除)
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"0"}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The mailbox cannot be empty",
		})
		return
	}
	// _, err := models.RDB.Get(c, email).Result()
	ttlResult, err := models.RDB.TTL(c, email).Result()
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The redis operation failed",
		})
		return
	} else if ttlResult == time.Duration(-1) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "No expiration time is set for the current Key",
		})
		return
	} else if ttlResult != time.Duration(-2) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The verification code has not expired, please use the previous one",
			"data": gin.H{
				"ttl": ttlResult.Seconds(),
			},
		})
		return
	}
	code := helper.GetRand()
	err = helper.SendCode(email, code)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Failed to send the verification code",
		})
		return
	}
	err = models.RDB.Set(c, email, code,
		time.Duration(config.Config.VerificationCodeDuration*int(time.Minute))).Err()
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Unable to write data to redis",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "The verification code is sent successfully",
		"data": gin.H{
			"ttl": time.Duration(config.Config.VerificationCodeDuration * int(time.Minute)).Seconds(),
		},
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param loginType formData string true "loginType"
// @Param nameOrMail formData string true "nameOrMail"
// @Param authCode formData string true "authCode"
// @Success 200 {string} json "{"code":"0"}"
// @Router /login [post]
func Login(c *gin.Context) {
	loginType := c.PostForm("loginType")
	nameOrMail := c.PostForm("nameOrMail")
	authCode := c.PostForm("authCode")
	if helper.IsNuiStrs(loginType, nameOrMail, authCode) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	user := &models.Users{}
	// 有三种登录方式：0 : 用户名+密码，1 : 邮箱+密码，2 : 邮箱+验证码
	if loginType == "0" || loginType == "1" {
		condition := "name"
		if loginType == "1" {
			condition = "email"
		}
		condition += " = ? AND password = ?"
		err := models.DB.Model(&models.Users{}).
			Where(condition, nameOrMail, authCode).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The user name does not exist or the password is incorrect",
				})
				return
			}
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mysql operation failure",
			})
			return
		}
	} else {
		sysCode, err := models.RDB.Get(c, nameOrMail).Result()
		if err != nil {
			if err == redis.Nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "The mailbox does not exist or the verification code has expired",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Redis operation failure",
				})
			}
			return
		}
		if sysCode != authCode {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code is incorrect",
			})
			return
		}
		err = models.DB.Model(&models.Users{}).
			Where("email = ?", nameOrMail).First(user).Error
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Mailbox does not exist",
			})
			return
		}
	}
	// 检查一下用户被封禁没有
	if !models.JudgePermit(models.Login, user.Ban) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The user has been banned",
		})
		return
	}
	// 上面是局部鉴权，下面是全局鉴权
	safeBan, err := models.RDB.Get(c, "safeBan").Result()
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Redis operation failed",
		})
		return
	}
	if safeBan != "permit" && !models.JudgePermit(models.Admin, user.Ban) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The website is currently in secure mode, and only administrators can log in",
		})
		return
	}
	token, err := middlewares.GenerateToken(user.ID)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Generate token fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Login success!",
		"data": gin.H{
			"token":    token,
			"name":     user.Name,
			"password": user.Password,
			"ban":      user.Ban, //前端根据用户权限，创建一些管理员特有的组件
		},
	})
}

// Register
// @Tags 公共方法
// @Summary 注册新用户，第一个注册的用户是管理员
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param registerTime formData string true "registerTime"
// @Success 200 {string} json "{"code":"0"}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	registerTime := c.PostForm("registerTime")
	if helper.IsNuiStrs(email, userCode, name, password, registerTime) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code has expired",
			})
		} else {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failure",
			})
		}
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The verification code is incorrect",
		})
		return
	}
	// 判断邮箱和用户名是否已存在
	var cnt int64
	err = models.DB.Model(&models.Users{}).
		Where("email = ? OR name = ?", email, name).Count(&cnt).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The email address or user name already exists",
		})
		return
	}
	var userCount int64
	err = models.DB.Model(&models.Users{}).Count(&userCount).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	userBan := 1
	if userCount == 0 {
		// 将第一位注册的用户升级为超级管理员
		userBan = 0
	}
	user := &models.Users{
		Name:         name,
		Email:        email,
		Password:     password,
		Score:        config.Config.UserInitScore,
		RegisterTime: registerTime,
		Ban:          userBan,
	}
	err = models.DB.Create(user).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	// 成功创建用户后，应该立即将验证码销毁掉，以免一把钥匙打开多道们的情况出现
	err = models.RDB.Unlink(c, email).Err()
	if err != nil {
		logger.Errorln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Account registration successful",
	})
}

// FindPassword
// @Tags 公共方法
// @Summary 找回密码
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"0"}"
// @Router /find-password [post]
func FindPassword(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	password := c.PostForm("password")
	if helper.IsNuiStrs(email, userCode, password) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Incomplete information",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The verification code has expired",
			})
		} else {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Redis operation failure",
			})
		}
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The verification code is incorrect",
		})
		return
	}
	user := &models.Users{}
	err = models.DB.Model(&models.Users{}).Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Email not registered",
			})
			return
		}
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	err = models.DB.Model(&models.Users{}).Where(&models.Users{ID: user.ID}).
		Update("password", password).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "Mysql operation failure",
		})
		return
	}
	err = models.RDB.Unlink(c, email).Err()
	if err != nil {
		logger.Errorln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Password changed successfully",
	})
}
