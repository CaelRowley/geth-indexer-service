package pubsub

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer(url string) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": url})
}

func (k *KafkaPubSub) StartProducerEventLoop() {
	for e := range k.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				slog.Error("producer delivery failed", "err", ev.TopicPartition.Error)
			}
		case kafka.Error:
			slog.Error("producer error", "code", ev.Code(), "err", ev.Error())
		}
	}
	slog.Info("kafka producer stopped")
}

func (k *KafkaPubSub) PublishBlock(blockData []byte) error {
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &blocksTopic, Partition: kafka.PartitionAny},
		Value:          blockData,
	}, nil)
	return err
}
