package main

import (
	kron "kron/internal"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var opsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "kron_failed_ops",
	Help: "The total number of failed operations",
}, []string{"consumer_id"})

var opsSuccessful = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "kron_successful_ops",
	Help: "The total number of successful operations",
}, []string{"consumer_id"})

var opsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "kron_ops_duration",
	Help:    "The duration of operations in milliseconds",
	Buckets: prometheus.DefBuckets,
}, []string{"consumer_id"})

var opsDelay = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "kron_ops_delay",
	Help:    "The delay of operations in milliseconds",
	Buckets: prometheus.DefBuckets,
}, []string{"consumer_id"})

func main() {
	KAFKA_BROKERS := os.Getenv("KAFKA_BROKERS")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")
	CONSUMER_ID := os.Getenv("CONSUMER_ID")

	logger := log.New(os.Stdout, "[Consumer] ", log.LstdFlags)

	// Kafka client
	brokers := strings.Split(KAFKA_BROKERS, ",")
	kafkaClient := kron.NewKafkaClient(brokers, KAFKA_TOPIC)

	// Prometheus metrics label
	label := prometheus.Labels{
		"consumer_id": CONSUMER_ID,
	}

	go func() {
		for {
			job, err := kafkaClient.ReadJob()
			if err != nil {
				logger.Println("Error reading message:", err)
				continue
			}
			logger.Println("Executing job:", job)
			startTime := time.Now()
			err = exec.Command("sh", "-c", job.Command).Run()
			endTime := time.Now()
			opsDuration.With(label).Observe(float64(endTime.Sub(startTime).Milliseconds()))
			opsDelay.With(label).Observe(float64(startTime.Sub(job.Time).Milliseconds()))
			if err != nil {
				logger.Println("Error executing command:", err)
				opsFailed.With(label).Inc()
				continue
			}
			opsSuccessful.With(label).Inc()
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
