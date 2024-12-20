package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Task struct {
	Schedule string
	Action   func() error
}

type Service struct {
	Tasks []Task
}

func (s *Service) Cron() {
	c := cron.New()

	for _, task := range s.Tasks {
		_, err := c.AddFunc(task.Schedule, func() {
			err := task.Action()
			if err != nil {
				log.Printf("Task failed: %v", err)
			} else {
				log.Println("Task completed successfully")
			}
		})
		if err != nil {
			log.Fatalf("Failed to add cron job: %v", err)
		}
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
