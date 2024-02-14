package helper

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/lazyliqiquan/help_rookie/config"
)

var (
	TokenPrivateKey []byte //token加密私钥
)

// 用来生成token的结构，主要用来鉴权
type UserClaims struct {
	Id  int
	Ban uint64
	jwt.StandardClaims
}

func init() {
	// 生成一个随机字符串，用来当作私钥，这样每次重启后的私钥都是不一样的，
	// 这样也会导致之前的token全部失效
	TokenPrivateKey = []byte(uuid.New().String())
}

// 生成 token
func GenerateToken(id int, ban uint64) (string, error) {
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

// 将给定字符串进行md5加密
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 生成唯一码
func GetUUID() string {
	return uuid.New().String()
}

// 生成验证码
func GetRand() string {
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// 向指定邮箱发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = config.Config.SenderMailbox
	e.To = []string{toUserEmail}
	e.Subject = "You received a verification code from help-cookie !"
	e.HTML = []byte("<b>" + code + "</b>")
	return e.Send(config.Config.SmtpServerPath+":"+config.Config.SmtpServerPort,
		smtp.PlainAuth("", config.Config.SenderMailbox,
			config.Config.SmtpServerVerification, config.Config.SmtpServerPath))
	// return e.SendWithTLS(config.Config.SmtpServerPath+":"+config.Config.SmtpServerPort,
	// 	smtp.PlainAuth("", config.Config.SenderMailbox,
	// 		config.Config.SmtpServerVerification, config.Config.SmtpServerPath),
	// 	&tls.Config{InsecureSkipVerify: true, ServerName: config.Config.SmtpServerPath})
}

// 读取文件为[]byte类型
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	// 读取文件内容
	return io.ReadAll(file)
}

func IsNuiStrs(arr ...string) bool {
	for _, e := range arr {
		if e == "" {
			return true
		}
	}
	return false
}
