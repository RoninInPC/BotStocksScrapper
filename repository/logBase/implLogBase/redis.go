package implLogBase

import (
	"BotStocksScrapper/app/entity"
	"BotStocksScrapper/hash"
	"BotStocksScrapper/repository/logBase"
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config entity.RedisDBConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}

type RedisRepository struct {
	client *redis.Client
	hash   hash.Hash
}

func NewRedisRepository(client *redis.Client, h hash.Hash) logBase.RedisRepository {
	return &RedisRepository{client: client, hash: h}
}

func (r *RedisRepository) Add(value string) bool {
	ctx := context.Background()
	key := r.hash(value)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	if exists == 1 {
		return false
	}
	err = r.client.Set(ctx, key, "value", 0).Err()
	return err == nil
}

func (r *RedisRepository) Free() bool {
	ctx := context.Background()
	err := r.client.FlushDB(ctx).Err()
	return err == nil
}
