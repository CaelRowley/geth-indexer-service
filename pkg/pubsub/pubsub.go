package pubsub

import (
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	blocksTopic = "blocks"
)

type PubSub interface {
	PublishBlock([]byte) error
	ConsumeBlock(*kafka.Message) error
	StartConsumer() error
}

type KafkaPubSub struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
	dbConn   db.DB
}

func NewPubSub(url string, dbConn db.DB) (PubSub, error) {
	p, err := NewProducer(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}
	c, err := NewConsumer(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	return &KafkaPubSub{p, c, dbConn}, nil
}
