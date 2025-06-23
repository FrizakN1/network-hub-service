package kafka

import (
	"backend/proto/searchpb"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type NodeProducer interface {
	SendBatchNodes(ctx context.Context, nodes []*searchpb.Node) error
	SendSingleNode(ctx context.Context, node *searchpb.Node) error
}

type DefaultNodeProducer struct {
	writer *kafka.Writer
}

type IndexNodeMessage struct {
	Type  string           `json:"type"`
	Node  *searchpb.Node   `json:"node"`
	Nodes []*searchpb.Node `json:"nodes"`
}

func NewNodeProducer(writer *kafka.Writer) NodeProducer {
	return &DefaultNodeProducer{
		writer: writer,
	}
}

func (p *DefaultNodeProducer) SendSingleNode(ctx context.Context, node *searchpb.Node) error {
	msg := &IndexNodeMessage{
		Type: "single",
		Node: node,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("single"),
		Value: data,
	})
}

func (p *DefaultNodeProducer) SendBatchNodes(ctx context.Context, nodes []*searchpb.Node) error {
	msg := &IndexNodeMessage{
		Type:  "batch",
		Nodes: nodes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("batch"),
		Value: data,
	})
}
