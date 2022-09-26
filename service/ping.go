package service

import "github.com/gin-gonic/gin"

func PingHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "pong"})
}
