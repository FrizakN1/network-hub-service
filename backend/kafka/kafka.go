package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"os"
	"strconv"
)

func CreateTopics() error {
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("KAFKA_ADDRESS"), os.Getenv("KAFKA_PORT")))
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

	topicConfigs := make([]kafka.TopicConfig, 2)

	topicConfigs[0] = kafka.TopicConfig{
		Topic:             "index-node",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	topicConfigs[1] = kafka.TopicConfig{
		Topic:             "index-hardware",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = controllerConn.CreateTopics(topicConfigs[0], topicConfigs[1])
	if err != nil {
		return fmt.Errorf("CreateTopics error (might be topic exists): %v\n", err)
	}

	return nil
}
