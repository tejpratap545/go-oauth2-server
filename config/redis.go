package config

import (
	"time"

	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

func Redis() *redis.Database {

	db := redis.New(redis.Config{
		Network:   "tcp",
		Addr:      "127.0.0.1:6379",
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  "",
		Database:  "",
		Prefix:    "",
	})

	return db
}
