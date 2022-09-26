package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var secretKey = []byte("123")

// 生成token
func TestGenerateToken(t *testing.T) {
	now := time.Now()
	expiresTime := now.Add(time.Hour * 3)
	claims := UserClaims{
		Name:     "user_1",
		Identity: "Get",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("tokenStr: %v\n", tokenStr)
}

// 解析token
func TestParseToken(t *testing.T) {
	s := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcl8xIiwiaWRlbnRpdHkiOiJHZXQifQ.xHe6R5E2Qii4zMvny8JZWhLCBWqnyuKHK-IwnxMZdTs"
	s = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcl8xIiwiaWRlbnRpdHkiOiJHZXQiLCJpc19hZG1pbiI6MCwiZXhwIjoxNjY0MTMxODQyfQ.eWQXqRxuxNFQ9RKCrarxhMhg5jmSrAcHqeug-n4k27M"
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(s, userClaims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Printf("userClaims: %+v\n", userClaims)
	}

}
