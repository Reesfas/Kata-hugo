package cache

import (
	"github.com/go-redis/redis"
)

func NewRedisClient(host, port string) *redis.Client {
	options := &redis.Options{
		Addr: host + ":" + port,
	}

	client := redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
