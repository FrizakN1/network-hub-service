package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
)

func NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%s", os.Getenv("KAFKA_ADDRESS"), os.Getenv("KAFKA_PORT"))),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
