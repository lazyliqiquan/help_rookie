package middlewares

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lazyliqiquan/help_rookie/config"
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
	Id int //只有用户id是不会变的
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

// 好像每次程序重新启动，之前的token全部失效了
// 生成 token
func GenerateToken(id int) (string, error) {
	UserClaim := &UserClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.Config.TokenDuration * int64(time.Hour))).Unix(),
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
