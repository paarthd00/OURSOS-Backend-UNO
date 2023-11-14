package redis

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"oursos.com/packages/util"
)

func Client() *redis.Client {
	err := godotenv.Load()

	util.CheckError(err)

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       0,
	})
	return client
}
