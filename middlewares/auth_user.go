package middlewares

import (
	"gin_gorm_oj/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		claims, err := helper.ParseToken(auth)
		if err != nil {
			ctx.Abort()
			log.Printf("parse token err:%v\n", err.Error())
			ctx.JSON(200, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized authorization",
			})
			return
		}
		if claims == nil {
			ctx.Abort()
			ctx.JSON(200, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized user",
			})
			return
		}
		ctx.Set("user", claims)
		ctx.Next()
	}
}
