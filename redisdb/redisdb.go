package redisdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	ctx          = context.Background()
	RedisClient  *redis.Client
	redisAddress string
)

func InitRedis() *redis.Client {
	const redisEnv = "REDIS_ADDR"
	redisAddress = os.Getenv(redisEnv)
	if redisAddress == "" {
		_, _ = fmt.Fprintln(os.Stderr, redisEnv)
		os.Exit(1)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	return client
}
