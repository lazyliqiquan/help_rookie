package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/helper"
	"github.com/lazyliqiquan/help_rookie/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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
			"msg":  "必填信息为空",
		})
		return
	}
	user := &models.Users{}
	// 有三种登录方式：0 : 用户名+密码，1 : 邮箱+密码，2 : 邮箱+验证码
	if loginType == "0" || loginType == "1" {
		authCode = helper.GetMd5(authCode)
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
					"code": 2,
					"msg":  "用户名不存在或密码错误",
				})
				return
			}
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "mysql操作失败",
			})
			return
		}
	} else {
		sysCode, err := models.RDB.Get(c, nameOrMail).Result()
		if err != nil {
			if err == redis.Nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 3,
					"msg":  "邮箱不存在或验证码过期",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 4,
					"msg":  "redis操作失败",
				})
			}
			return
		} else if sysCode != authCode {
			c.JSON(http.StatusOK, gin.H{
				"code": 5,
				"msg":  "验证码不正确",
			})
			return
		}
		err = models.DB.Model(&models.Users{}).
			Where("email = ?", nameOrMail).First(user).Error
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 6,
				"msg":  "邮箱不存在",
			})
			return
		}
	}
	token, err := helper.GenerateToken(user.ID, user.Ban)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 7,
			"msg":  "Generate token fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登陆成功!",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"0"}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "邮箱不能为空",
		})
		return
	}
	_, err := models.RDB.Get(c, email).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  "请勿重复请求验证码，请使用之前的验证码",
		})
		return
	} else if err != redis.Nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  "获取redis数据失败",
		})
		return
	}
	code := helper.GetRand()
	err = helper.SendCode(email, code)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4,
			"msg":  "发送验证码失败",
		})
		return
	}
	err = models.RDB.Set(c, email, code,
		time.Duration(config.Config.VerificationCodeDuration*int(time.Minute))).Err()
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 5,
			"msg":  "无法向redis写入数据",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "验证码发送成功",
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
			"msg":  "参数不正确",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "验证码已过期",
			})
		} else {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 3,
				"msg":  "redis异常",
			})
		}
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 4,
			"msg":  "验证码不正确",
		})
		return
	}
	// 判断邮箱是否已存在
	var cnt int64
	err = models.DB.Model(&models.Users{}).Where("email = ? OR name = ?", email, name).Count(&cnt).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 5,
			"msg":  "mysql异常",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 6,
			"msg":  "该邮箱已被注册",
		})
		return
	}
	var userCount int64
	err = models.DB.Model(&models.Users{}).Count(&userCount).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 5,
			"msg":  "mysql异常",
		})
		return
	}
	userBan := config.Config.UserBan
	if userCount == 0 {
		// 将第一位注册的用户升级为超级管理员
		userBan |= 1
	}
	user := &models.Users{
		Name:         name,
		Email:        email,
		Password:     helper.GetMd5(password),
		Score:        config.Config.UserInitScore,
		RegisterTime: registerTime,
		Ban:          userBan,
	}
	err = models.DB.Create(user).Error
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 5,
			"msg":  "mysql异常",
		})
		return
	}
	token, err := helper.GenerateToken(user.ID, user.Ban)
	if err != nil {
		logger.Errorln(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 6,
			"msg":  "Token生成失败",
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
		"msg":  "账号注册成功",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
