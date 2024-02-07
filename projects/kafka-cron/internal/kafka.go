package kron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	brokers []string
	topic   string
	r       *kafka.Reader
	w       *kafka.Writer
}

func NewKafkaClient(brokers []string, topic string) *KafkaClient {
	log.Println("Connecting to Kafka brokers:", brokers)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	k := &KafkaClient{
		brokers: brokers,
		topic:   topic,
		r:       r,
		w:       w,
	}
	return k
}

func (k *KafkaClient) CreateTopic() error {
	log.Println("Dialing Kafka broker:", k.brokers[0])
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             k.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}

func (k *KafkaClient) WriteJob(job Job) error {
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Value: b,
	}
	return k.w.WriteMessages(context.Background(), message)
}

func (k *KafkaClient) ReadJob() (*Job, error) {
	m, err := k.r.ReadMessage(context.Background())
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reader closed, topic: %s", k.r.Config().Topic)
		}
		return nil, err
	}
	var job Job
	err = json.Unmarshal(m.Value, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (k *KafkaClient) Close() error {
	werr := k.w.Close()
	rerr := k.r.Close()
	if werr != nil {
		return werr
	}
	return rerr
}
