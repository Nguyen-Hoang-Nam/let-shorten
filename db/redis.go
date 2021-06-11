package db

import (
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func RedisInit() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetRedis() *redis.Client {
	return client
}
