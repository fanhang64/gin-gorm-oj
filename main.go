package main

import (
	"gin_gorm_oj/models"
	"gin_gorm_oj/router"
	"log"
)

func main() {
	r := router.SetupRouter()
	err := models.InitDB()
	if err != nil {
		log.Printf("init db error:%v\n", err.Error())
		return
	}
	models.InitRedisDB()

	r.Run()
}
