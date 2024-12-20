package usecase

import (
	"errors"
	"log"
	"reddis/app/repo/changeBase"
	"reddis/app/repo/logBase"
)

type Cleaner interface {
	Clean() error
}

type DatabaseCleaner struct {
	logBaseRepo    logBase.RedisRepository
	changeBaseRepo changeBase.CBRepository
}

func NewDatabaseCleaner(logBaseRepo logBase.RedisRepository, changeBaseRepo changeBase.CBRepository) *DatabaseCleaner {
	return &DatabaseCleaner{
		logBaseRepo:    logBaseRepo,
		changeBaseRepo: changeBaseRepo,
	}
}

func (c *DatabaseCleaner) Clean() error {
	if !c.logBaseRepo.Free() {
		return errors.New("failed to clean LogBase")
	} else {
		log.Println("Clean LogBase")
	}

	if !c.changeBaseRepo.Free() {
		return errors.New("failed to clean ChangeBase")
	} else {
		log.Println("Clean ChangeBase")
	}

	return nil
}
