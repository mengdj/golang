package model

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"tool"
)


//JWT验证
type User struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	jwt.StandardClaims
}

//获取ken
func GetToken(usr, pwd string) (string, error) {
	ins := User{Account: usr, Password: pwd}
	ins.ExpiresAt = time.Now().Add(tool.TOKENEXPIREDURATION).Unix()
	ins.Issuer = "mengdj"
	return jwt.NewWithClaims(jwt.SigningMethodHS256, ins).SignedString([]byte(tool.TOKENSECRECT))
}

//解析token
func ParseToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(tool.TOKENSECRECT), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*User); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
