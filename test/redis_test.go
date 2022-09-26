package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

var redisDB *redis.Client
var ctx = context.Background()

func TestConnectRedis(t *testing.T) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	result, err := redisDB.Ping(ctx).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %v\n", result)

	err = redisDB.Set(ctx, "name", "zs", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	result, err = redisDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %v\n", result)
}
