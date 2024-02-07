package main

import (
	kron "kron/internal"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	KAFKA_BROKERS := os.Getenv("KAFKA_BROKERS")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")

	logger := log.New(os.Stdout, "[Consumer] ", log.LstdFlags)

	// Kafka client
	brokers := strings.Split(KAFKA_BROKERS, ",")
	kafkaClient := kron.NewKafkaClient(brokers, KAFKA_TOPIC)

	for {
		job, err := kafkaClient.ReadJob()
		if err != nil {
			logger.Println("Error reading message:", err)
			continue
		}
		logger.Println("Executing job:", job)
		err = exec.Command("sh", "-c", job.Command).Run()
		if err != nil {
			logger.Println("Error executing command:", err)
		}
	}
}
