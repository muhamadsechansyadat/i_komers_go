package config

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func SetupRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatalf("Error converting Redis DB to integer: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       redisDB,
	})

	if err != nil {
		panic(err.Error())
	}

	return client
}
