package handlers

import (
	"backend/database"
	"backend/kafka"
	"backend/proto/searchpb"
	"context"
)

type NodeLogic interface {
	SendBatchNodes(ctx context.Context) error
}

type DefaultNodeLogic struct {
	NodeRepo database.NodeRepository
	kafka.NodeProducer
}

func NewNodeLogic(db *database.Database) NodeLogic {
	return &DefaultNodeLogic{
		NodeRepo: &database.DefaultNodeRepository{
			Database: *db,
		},
		NodeProducer: kafka.NewNodeProducer(kafka.NewKafkaWriter("index-node")),
	}
}

func (l *DefaultNodeLogic) SendBatchNodes(ctx context.Context) error {
	nodes, err := l.NodeRepo.GetNodesForIndex()
	if err != nil {
		return err
	}

	var grpcNodes []*searchpb.Node

	for _, node := range nodes {
		grpcNode := &searchpb.Node{
			Id:    int32(node.ID),
			Name:  node.Name,
			Zone:  node.Zone.String,
			Owner: node.Owner.Value,
		}

		grpcNodes = append(grpcNodes, grpcNode)
	}

	const batchSize = 1000

	for i := 0; i < len(grpcNodes); i += batchSize {
		end := i + batchSize
		if end > len(grpcNodes) {
			end = len(grpcNodes)
		}

		batch := grpcNodes[i:end]

		if err = l.NodeProducer.SendBatchNodes(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}
