package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reddis/app/controller/cron"
	"reddis/app/entity"
	"reddis/app/repo/implementation"
)

func main() {

	config := loadConfig("config/config.yml")

	redisClient := implementation.NewRedisClient(config.Redis)

	redisRepo := implementation.NewRedisRepository(redisClient)

	if redisRepo.Add("mykey") {
		log.Println("String added successfully")
	} else {
		log.Println("String already exists or failed to add")
	}

	if redisRepo.Free() {
		log.Println("Database cleared successfully")
	} else {
		log.Println("Failed to clear database")
	}

	cleaner := cron.NewRedisCleaner(redisRepo)

	cron.StartCronService(cleaner)
}

func loadConfig(filename string) *entity.Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config entity.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &config
}
