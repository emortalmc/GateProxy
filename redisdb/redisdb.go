package redisdb

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	ctx         = context.Background()
	RedisClient *redis.Client
)

func InitRedis() *redis.Client {
	//addr := flag.String("redisaddress", "localhost:6379", "The address of Redis server")
	//flag.Parse()
	client := redis.NewClient(&redis.Options{
		Addr: "172.17.0.1:6379",
		//Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	return client
}
