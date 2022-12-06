package redisdb

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
)

var (
	addr        = flag.String("addr", "localhost:6379", "The address of Redis server")
	ctx         = context.Background()
	RedisClient = initRedis()
)

func initRedis() *redis.Client {
	flag.Parse()
	client := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: "",
		DB:       0,
	})

	return client
}
