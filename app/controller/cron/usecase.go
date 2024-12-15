package cron

import (
	"errors"
	"reddis/app/repo"
)


type RedisCleaner struct {
	repo repo.RedisRepository
}

func NewRedisCleaner(repo repo.RedisRepository) *RedisCleaner {
	return &RedisCleaner{repo: repo}
}

func (c *RedisCleaner) Clean() error {
	if !c.repo.Free() {
		return errors.New("failed to clean redis")
	}

	return nil
}
