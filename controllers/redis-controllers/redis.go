package redis_controllers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitializeRedis(ctx context.Context) {
	Rdb = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
	})

	// Ping the Redis server to check if it's reachable
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
			panic("failed to connect Redis")
	}
	fmt.Println("Redis connected successfully")
}