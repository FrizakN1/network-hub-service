package kafka

import (
	"backend/proto/searchpb"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type HardwareProducer interface {
	SendBatchHardware(ctx context.Context, hardware []*searchpb.Hardware) error
	SendSingleHardware(ctx context.Context, hardware *searchpb.Hardware) error
}

type DefaultHardwareProducer struct {
	writer *kafka.Writer
}

type IndexHardwareMessage struct {
	Type           string               `json:"type"`
	HardwareSingle *searchpb.Hardware   `json:"hardware_single"`
	Hardware       []*searchpb.Hardware `json:"hardware"`
}

func NewHardwareProducer(writer *kafka.Writer) HardwareProducer {
	return &DefaultHardwareProducer{
		writer: writer,
	}
}

func (p *DefaultHardwareProducer) SendSingleHardware(ctx context.Context, hardware *searchpb.Hardware) error {
	msg := &IndexHardwareMessage{
		Type:           "single",
		HardwareSingle: hardware,
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

func (p *DefaultHardwareProducer) SendBatchHardware(ctx context.Context, hardware []*searchpb.Hardware) error {
	msg := &IndexHardwareMessage{
		Type:     "batch",
		Hardware: hardware,
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
