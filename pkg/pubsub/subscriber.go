//go:build kafka

package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Subscriber interface {
	StartPoll(context.Context) error
	Close() error
}

type KafkaConsumer struct {
	*kafka.Consumer
	dbConn db.DB
}

func NewSubscriber(url string, dbConn db.DB) (Subscriber, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        url,
		"group.id":                 "evmIndexer",
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{blocksTopic, txsTopic}, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to kafka topics: %w", err)
	}

	return &KafkaConsumer{c, dbConn}, nil
}

func (c *KafkaConsumer) StartPoll(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("kafka consumer stopped")
			return nil
		default:
			ev := c.Consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if *e.TopicPartition.Topic == blocksTopic {
					if err := c.handleBlock(e); err != nil {
						slog.Error("failed to consume block message", "err", err)
					}
				}
				if *e.TopicPartition.Topic == txsTopic {
					if err := c.handleTx(e); err != nil {
						slog.Error("failed to consume tx message", "err", err)
					}
				}
			case kafka.Error:
				fmt.Println(e)
				slog.Error("kafka subscription failed", "code", e.Code(), "err", e.Error())
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			}
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.Consumer.Close()
}

func (c *KafkaConsumer) handleBlock(m *kafka.Message) error {
	var block data.Block
	if err := json.Unmarshal(m.Value, &block); err != nil {
		return fmt.Errorf("failed to unmarshal block data: %w", err)
	}
	if err := c.dbConn.InsertBlock(block); err != nil {
		return fmt.Errorf("failed to store block in db: %w", err)
	}
	if _, err := c.Consumer.StoreMessage(m); err != nil {
		return fmt.Errorf("failed to store kafka offset after message: %w", err)
	}
	return nil
}

func (c *KafkaConsumer) handleTx(m *kafka.Message) error {
	var tx data.Transaction
	if err := json.Unmarshal(m.Value, &tx); err != nil {
		return fmt.Errorf("failed to unmarshal tx data: %w", err)
	}
	if err := c.dbConn.InsertTx(tx); err != nil {
		return fmt.Errorf("failed to store tx in db: %w", err)
	}
	if _, err := c.Consumer.StoreMessage(m); err != nil {
		return fmt.Errorf("failed to store kafka offset after message: %w", err)
	}
	return nil
}