package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"telegram-bot-golang/env"
)

var ctx = context.Background()

func getRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Get(key string) (value string, err error) {
	if env.CacheIsEnabled() {
		value, err = getRedis().
			Get(ctx, key).
			Result()
		if err != nil {
			fmt.Println("redis:error of getting cache:" + err.Error())
		}
		return
	}
	err = redis.ErrClosed
	return
}

func Set(key string, value interface{}) {
	err := getRedis().Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("redis:error of writing cache:" + err.Error())
	}
}

func Del(key string) {
	err := getRedis().Del(ctx, key).Err()
	if err != nil {
		fmt.Println("redis:error of deleting cache:" + err.Error())
	}
}
