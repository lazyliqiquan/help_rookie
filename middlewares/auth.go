package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/models"
	"go.uber.org/zap"
)

type AuthType int

const (
	UserLevel AuthType = iota
	AdminLevel
	RootLevel
)

var (
	TokenPrivateKey []byte //token加密私钥
	logger          *zap.SugaredLogger
)

// 用来生成token的结构，主要用来鉴权
type UserClaims struct {
	Id  int
	Ban int
	jwt.StandardClaims
}

func Init(loggerInstance *zap.SugaredLogger) {
	logger = loggerInstance
}

func init() {
	// 生成一个随机字符串，用来当作私钥，这样每次重启后的私钥都是不一样的，
	// 这样也会导致之前的token全部失效
	TokenPrivateKey = []byte(uuid.New().String())
}

func AuthCheck(authType AuthType) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Unauthorized",
			})
			return
		}
		userClaim, err := AnalyseToken(auth)
		if err != nil {
			logger.Errorln(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "Operation failed, please log in first",
			})
			return
		}
		if !models.IsLogin(userClaim.Ban) {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "The user has been banned",
			})
			return
		}
		if authType == AdminLevel {
			if !models.IsAdmin(userClaim.Ban) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Operation failed. Administrator rights required",
				})
				return
			}
		} else if authType == RootLevel {
			if !models.IsRoot(userClaim.Ban) {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "Operation failed. Root rights required",
				})
				return
			}
		}
		c.Set("id", userClaim.Id)
		c.Set("ban", userClaim.Ban)
		c.Next()
	}
}

// 生成 token
func GenerateToken(id, ban int) (string, error) {
	UserClaim := &UserClaims{
		Id:  id,
		Ban: ban,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: config.Config.TokenDuration * int64(time.Hour),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(TokenPrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim,
		func(token *jwt.Token) (interface{}, error) {
			return TokenPrivateKey, nil
		})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("Analyse token fail")
	}
	return userClaim, nil
}
