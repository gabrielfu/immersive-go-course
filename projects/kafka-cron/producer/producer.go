package main

import (
	"log"
	"os"
	"strings"
)

const (
	BUFFER_SIZE = 10
	PORT        = "8080"
)

func main() {
	KAFKA_BROKERS := os.Getenv("KAFKA_BROKERS")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")

	brokers := strings.Split(KAFKA_BROKERS, ",")
	producer := NewKafkaClient(brokers, KAFKA_TOPIC)

	// Job scheduler
	scheduler := NewScheduler(BUFFER_SIZE)
	scheduler.Start()
	defer scheduler.Stop()

	// When the scheduler fires, produce message to Kafka
	go func() {
		for job := range scheduler.JobsDue {
			log.Println("Producing message to Kafka:", job)
			producer.WriteJob(job)
		}
	}()

	var jobHandler = func(schedule, command string) error {
		s, err := ParseCronSchedule(schedule)
		if err != nil {
			return err
		}
		return scheduler.Add(s, command)
	}

	ServeAPI(jobHandler, PORT)
}
