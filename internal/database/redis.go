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
		Addr:             config.RedisURL,
		PoolSize:         100,                  
		MinIdleConns:     32,                   
		MaxRetries:       5,                    
		MinRetryBackoff:  10 * time.Millisecond, 
		MaxRetryBackoff:  512 * time.Millisecond,
		DialTimeout:      5 * time.Second,      
		ReadTimeout:      500 * time.Millisecond, 
		WriteTimeout:     500 * time.Millisecond, 
		PoolTimeout:      4 * time.Second,      
		ConnMaxIdleTime:  5 * time.Minute,      
		ConnMaxLifetime:  30 * time.Minute,     
	}
	RedisClient = redis.NewClient(options)
}
