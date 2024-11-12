package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"micro-todoList-k8s/config"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// InitCache 在中间件中初始化redis链接
func InitCache() {
	host := config.C.Redis.RedisHost
	port := config.C.Redis.RedisPort
	password := config.C.Redis.RedisPassword
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	RedisClient = client
}
