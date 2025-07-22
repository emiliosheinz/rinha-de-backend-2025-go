package database

import (
	"context"
	"time"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisContext = context.Background()

func InitRedis() {
	options := &redis.Options{
		Addr:         config.RedisURL,
		PoolSize:     100,
		MinIdleConns: 25,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	}
	RedisClient = redis.NewClient(options)
}
