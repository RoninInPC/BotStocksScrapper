package implementation

import (
	"context"
	"github.com/redis/go-redis/v9"
	"reddis/app/controller"
	"reddis/app/repo"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) repo.RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Add(value string) bool {
	ctx := context.Background()
	key, err := controller.Hash(value)
	if err != nil {
		return false
	}
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
