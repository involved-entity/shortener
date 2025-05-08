package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v9"
)

var redisClient *redis.Client

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: " + err.Error())
		os.Exit(1)
	}
}

func GetClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client not initialized")
		os.Exit(1)
	}
	return redisClient
}
