package helper

import (
	"crypto/md5"
	"fmt"
	"gin_gorm_oj/constants"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken
// 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	now := time.Now()
	expiresTime := now.Add(time.Hour * 3) // 3个小时过期
	claims := UserClaims{
		Name:     name,
		Identity: identity,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(constants.SecretKey)
}

// ParseToken
// 解析token
func ParseToken(token string) (*UserClaims, error) {
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return constants.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("token error:%v\n", err.Error())
	}
	return userClaims, nil
}

// GetUUID
// 生成唯一标识符
func GetUUID() string {
	return uuid.NewV4().String()
}

// GenerateVerifyCode 生成验证码
func GenerateVerifyCode() string {
	var verifyCode string
	for i := 0; i < 6; i++ {
		verifyCode += fmt.Sprintf("%d", rand.Intn(10))
	}
	return verifyCode
}
