package main

import (
	"log"
	"os"
	"strings"

	kron "kron/internal"
)

const (
	BUFFER_SIZE = 10
	PORT        = "8080"
)

func main() {
	KAFKA_BROKERS := os.Getenv("KAFKA_BROKERS")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")

	logger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	// Kafka producer
	brokers := strings.Split(KAFKA_BROKERS, ",")
	producer := kron.NewKafkaClient(brokers, KAFKA_TOPIC)

	// Job scheduler
	scheduler := kron.NewScheduler(BUFFER_SIZE)
	scheduler.Start()
	defer scheduler.Stop()

	// When the scheduler fires, produce message to Kafka
	go func() {
		for job := range scheduler.JobsDue {
			logger.Println("Scheduler triggered job:", job)
			producer.WriteJob(job)
		}
	}()

	var jobHandler = func(schedule, command string) error {
		s, err := kron.ParseCronSchedule(schedule)
		if err != nil {
			return err
		}
		return scheduler.Add(s, command)
	}

	kron.ServeAPI(jobHandler, PORT)
}
