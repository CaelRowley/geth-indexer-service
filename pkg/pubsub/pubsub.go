package pubsub

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	blocksTopic = "blocks"
)

type PubSub interface {
	PublishBlock([]byte) error
	ConsumeBlock(*kafka.Message) error
	StartProducerEventLoop()
	StartConsumerPoll(context.Context) error
	Close() error
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

func (k *KafkaPubSub) Close() error {
	i := k.Producer.Flush(10000)
	for i > 0 {
		slog.Info("flushing kafka messages...", "remaining:", i)
		i = k.Producer.Flush(10000)
	}
	k.Producer.Close()
	return k.Consumer.Close()
}
