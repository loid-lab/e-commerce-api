package initializers

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

var ErrMissingRedisUrl = errors.New("REDIS_URL is not set in environment")

var RedisCLient *redis.Client

func RedisConnect() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		return nil, ErrMissingRedisUrl
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	RedisCLient = client
	return client, nil
}
