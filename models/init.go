package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var RedisDB *redis.Client

func InitDB() (err error) {
	dsn := "root:123@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=true&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("gorm init err:%v", err.Error())
		return
	}
	return
}

func InitRedisDB() {
	ctx := context.Background()
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	result, err := RedisDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("init redis db err:%v\n", err.Error())
	}
	fmt.Printf("result: %v\n", result)
}
