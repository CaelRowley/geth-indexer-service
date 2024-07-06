package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewConsumer(url string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        url,
		"group.id":                 "evmIndexer",
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
}

func (k *KafkaPubSub) StartConsumerPoll(ctx context.Context) error {
	if err := k.Consumer.SubscribeTopics([]string{blocksTopic}, nil); err != nil {
		return fmt.Errorf("failed to subscribe to kafka topics: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			slog.Info("kafka consumer stopped")
			return nil
		default:
			ev := k.Consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if *e.TopicPartition.Topic == blocksTopic {
					if err := k.ConsumeBlock(e); err != nil {
						slog.Error("failed to consume block message", "err", err)
					}
				}
			case kafka.Error:
				slog.Error("kafka subscription failed", "code", e.Code(), "err", e.Error())
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			}
		}
	}
}

func (k *KafkaPubSub) ConsumeBlock(m *kafka.Message) error {
	var block data.Block
	if err := json.Unmarshal(m.Value, &block); err != nil {
		return fmt.Errorf("failed to unmarshal block data: %w", err)
	}
	if err := k.dbConn.InsertBlock(block); err != nil {
		return fmt.Errorf("failed to store block in db: %w", err)
	}
	if _, err := k.Consumer.StoreMessage(m); err != nil {
		return fmt.Errorf("failed to store kafka offset after message: %w", err)
	}
	return nil
}
