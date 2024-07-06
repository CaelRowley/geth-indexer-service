package pubsub

import (
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

func (k *KafkaPubSub) StartConsumer() error {
	err := k.Consumer.SubscribeTopics([]string{blocksTopic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to kafka topics: %w", err)
	}

	for {
		ev := k.Consumer.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			if e.TopicPartition.Topic == &blocksTopic {
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

func (k *KafkaPubSub) ConsumeBlock(m *kafka.Message) error {
	var block data.Block
	err := json.Unmarshal(m.Value, &block)
	if err != nil {
		return fmt.Errorf("failed to unmarshal block data: %w", err)
	}

	err = k.dbConn.InsertBlock(block)
	if err != nil {
		return fmt.Errorf("failed to store block in db: %w", err)
	}

	_, err = k.Consumer.StoreMessage(m)
	if err != nil {
		return fmt.Errorf("failed to store kafka offset after message: %w", err)
	}

	return nil
}
