package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reddis/app/controller/cron"
	"reddis/app/entity"
	"reddis/app/repo/changeBase/impl"
	"reddis/app/repo/logBase/implLogBase"
	"reddis/app/usecase"
)

func main() {
	config := loadConfig("config/config.yml")

	redisClient := implLogBase.NewRedisClient(config.Redis.LogBase)
	changeBaseClient := impl.NewChangeBaseClient(config.Redis.ChangeBase)

	logBaseRepo := implLogBase.NewRedisRepository(redisClient)
	changeBaseRepo := impl.NewChangeBaseRedisRepository(changeBaseClient)

	if logBaseRepo.Add("mykey") {
		log.Println("String added successfully to LogBase")
	} else {
		log.Println("String already exists or failed to add to LogBase")
	}

	stock := entity.StockAdd{
		StockName: "AAPL",
		Type:      "SALE",
		NumPrice:  100,
	}
	if changeBaseRepo.Add(stock) {
		log.Println("Stock info added successfully to ChangeBase")
	} else {
		log.Println("Failed to add stock info to ChangeBase")
	}

	sum := changeBaseRepo.Get("AAPL", "SALE")
	log.Printf("Total SALE for AAPL: %d\n", sum)

	cleaner := usecase.NewDatabaseCleaner(logBaseRepo, changeBaseRepo)

	tasks := []cron.Task{
		{
			Schedule: "0 3 * * *", // Каждый день в 3:00 утра
			Action:   cleaner.Clean,
		},
	}

	cronService := cron.Service{
		Tasks: tasks,
	}

	cronService.Cron()
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
