package drivers

import (
	"github.com/go-redis/redis"
	"github.com/khalil-farashiani/url-shortener/internal/utils"
)

var Client *redis.Client

func init() {
	//Initializing redis
	dsn := utils.GetEnv("REDIS_DSN", "localhost:6379")
	Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
