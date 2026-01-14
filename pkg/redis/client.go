package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var AllCarsRedis = "all_cars"

var AllCarsExpirationTime = time.Minute * 30

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
