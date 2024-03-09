package main

import (
	"context"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/clean"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-redis/redis"
)

// SomeRepository интерфейс для получения данных
type SomeRepository interface {
	GetData(request string) (string, error)
}

// SomeRepositoryImpl структура для получения данных из базы
type SomeRepositoryImpl struct {
	api *clean.Api
}

func NewRepositoryImpl(api *clean.Api) *SomeRepositoryImpl {
	return &SomeRepositoryImpl{api}
}

func (r *SomeRepositoryImpl) GetData(request string) (string, error) {
	addresses, err := r.api.Address(context.Background(), request)
	if err != nil {
		return "", err
	}
	fmt.Println(addresses[0].Result)
	return addresses[0].Result, err
}

// SomeRepositoryProxy прокси для кэширования данных
type SomeRepositoryProxy struct {
	repository SomeRepository
	cache      *redis.Client
}

func NewRepositoryProxy(repository SomeRepository, cache *redis.Client) *SomeRepositoryProxy {
	return &SomeRepositoryProxy{repository: repository, cache: cache}
}

func (r *SomeRepositoryProxy) GetData(request string) (string, error) {
	// Проверяем наличие данных в кэше
	cachedData, err := r.cache.Get(request).Result()
	if err == nil && cachedData != "" {
		return cachedData, err
	}
	// Если данные отсутствуют в кэше, получаем их из оригинального объекта
	data, err := r.repository.GetData(request)
	if err != nil {
		return "", err
	}
	// Сохраняем данные в кэше
	err = r.cache.Set(request, data, 0).Err()
	if err != nil {
		return "", err
	}

	return data, nil
}

func main() {
	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "11cb4969967b7e68ab87b57258372aefec0eb6ac",
		SecretKeyValue: "3461265109aaa28b20523e1b4dfb4d36e475fc9f"}))
	// Создаем клиент Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Подключаемся к Redis
	_, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	// Создаем прокси для репозитория
	repository := NewRepositoryImpl(cleanApi)
	proxy := NewRepositoryProxy(repository, redisClient)

	// Получаем данные через прокси
	data, err := proxy.GetData("Москва")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Data:", data)
}
