package initializers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrMissingRedisUrl = errors.New("REDIS_URL is not set in environment")

var RedisCLient *redis.Client

func RedisConnect() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // or from .env when needed
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return RedisClient, RedisClient.Ping(ctx).Err()
}
