package test

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type userClaims struct {
	Identity string
	Name     string
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

func TestGenerateToken(t *testing.T) {
	userClaim := &userClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: time.Now().Add(time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZGVudGl0eSI6InVzZXJfMSIsIk5hbWUiOiJHZXQiLCJleHAiOjE3MDc3OTQ0MzR9.dSHl5X6W55RpnoegGt6SuEPqzNgM6vaJlM8a9AAQTu8
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf(tokenString)
	fmt.Println(tokenString)
}

func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZGVudGl0eSI6InVzZXJfMSIsIk5hbWUiOiJHZXQifQ.O8q7NCOLZp9Bgk-qoPQ68eE5N7r5jzEZ1tvmXVly7u4"
	userClaim := new(userClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}
