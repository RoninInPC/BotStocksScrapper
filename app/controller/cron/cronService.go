package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartCronService(cleaner *RedisCleaner) {
	c := cron.New()

	// Добавление задачи в крон (например, каждый день в 3:00 утра)
	_, err := c.AddFunc("0 3 * * *", func() {
		err := cleaner.Clean()
		if err != nil {
			log.Printf("Failed to clean Redis: %v", err)
		} else {
			log.Println("Redis cleaned successfully")
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Cron service started in the background.")
	<-stop

	log.Println("Shutting down...")
	c.Stop()

	log.Println("Cron service stopped. Exiting.")
}
