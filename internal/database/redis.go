package database

import (
	"context"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisContext = context.Background()

func InitRedis() {
	options := &redis.Options{Addr: config.RedisURL}
	RedisClient = redis.NewClient(options)
}
