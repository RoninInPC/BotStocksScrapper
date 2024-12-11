package redis

import (
	"github.com/redis/go-redis/v9"
	"reddis/app/entity"
)

func NewRedisClient(config entity.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
