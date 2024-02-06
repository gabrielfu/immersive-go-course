package main

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	brokers []string
	topic   string
	w       *kafka.Writer
}

func NewKafkaClient(brokers []string, topic string) *KafkaClient {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	k := &KafkaClient{
		brokers: brokers,
		topic:   topic,
		w:       w,
	}
	k.createTopic()
	return k
}

func (k *KafkaClient) createTopic() error {
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

func (k *KafkaClient) Close() error {
	return k.w.Close()
}
