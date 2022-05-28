package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Get(key string) string {
	value, err := getRedis().Get(ctx, key).
		Result()
	if err != nil {
		fmt.Println("error of getting cache:" + err.Error())
	}
	return value
}

func Set(key string, value interface{}) {
	err := getRedis().Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("error of writing cache:" + err.Error())
	}
}
