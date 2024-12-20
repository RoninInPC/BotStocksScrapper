package impl

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"reddis/app/entity"
)

type ChangeBaseRedisRepository struct {
	client *redis.Client
}

func NewChangeBaseClient(config entity.RedisDBConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
func NewChangeBaseRedisRepository(client *redis.Client) *ChangeBaseRedisRepository {
	return &ChangeBaseRedisRepository{client: client}
}

func (r *ChangeBaseRedisRepository) Add(stock entity.StockAdd) bool {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", stock.StockName, stock.Type)

	err := r.client.IncrBy(ctx, key, stock.NumPrice).Err()
	return err == nil
}

func (r *ChangeBaseRedisRepository) Get(stockName, operationType string) int64 {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", stockName, operationType)

	value, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		return 0
	}
	return value
}

func (r *ChangeBaseRedisRepository) Free() bool {

	ctx := context.Background()
	err := r.client.FlushDB(ctx).Err()
	return err == nil
}
