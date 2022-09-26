package helper

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Name     string `json:"username"`
	Identity string `json:"identity"`

	IsAdmin int `json:"is_admin"` // 是否管理员
	jwt.StandardClaims
}
