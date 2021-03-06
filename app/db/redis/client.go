package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"telegram-bot-golang/env"
	"time"
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

func SetStruct(key string, value interface{}, expiration time.Duration) {
	bytes, err := json.Marshal(value)
	if err != nil {
		fmt.Println("redis:marshal:error:" + err.Error())
		return
	}
	err = getRedis().Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		fmt.Println("redis:error of writing cache:" + err.Error())
	}
}

func Set(key string, value string, expiration time.Duration) {
	err := getRedis().Set(ctx, key, value, expiration).Err()
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
