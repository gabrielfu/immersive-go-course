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

	logger := log.New(os.Stdout, "[Producer] ", log.LstdFlags)

	// Kafka client
	brokers := strings.Split(KAFKA_BROKERS, ",")
	kafkaClient := kron.NewKafkaClient(brokers, KAFKA_TOPIC)
	err := kafkaClient.CreateTopic()
	if err != nil {
		logger.Fatalln("Error creating Kafka topic:", err)
	}

	// Job scheduler
	scheduler := kron.NewScheduler(BUFFER_SIZE)
	scheduler.Start()
	defer scheduler.Stop()

	// When the scheduler fires, produce message to Kafka
	go func() {
		for job := range scheduler.JobsDue {
			logger.Println("Scheduler triggered job:", job)
			err := kafkaClient.WriteJob(job)
			if err != nil {
				logger.Println("Error writing job to Kafka:", err)
				continue
			}
			logger.Println("Job written to Kafka:", job)
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
