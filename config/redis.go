package config

import (
	"github.com/go-redis/redis/v8"
)

func Redis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
