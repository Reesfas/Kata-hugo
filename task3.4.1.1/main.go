package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Cacher interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type cachestr struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cacher {
	return &cachestr{
		client: client,
	}
}

func (c *cachestr) Set(key string, value interface{}) error {
	err := c.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cachestr) Get(key string) (interface{}, error) {
	val, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("not found by key %s", key)
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func main() {
	// Создание клиента Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cache := NewCache(client)

	// Установка значения по ключу
	err := cache.Set("some:key", "value")
	if err != nil {
		panic(err)
	}

	// Получение значения по ключу
	value, err := cache.Get("some:key")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

}
