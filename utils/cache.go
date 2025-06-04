package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func SetJSON(client *redis.Client, key string, data any, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return client.Set(Ctx, key, bytes, ttl).Err()
}

func GetJSON[T any](client *redis.Client, key string, dest *T) error {
	val, err := client.Get(Ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func InvalidateKeys(client *redis.Client, key ...string) error {
	ctx := context.Background()
	return client.Del(ctx, key...).Err()
}
