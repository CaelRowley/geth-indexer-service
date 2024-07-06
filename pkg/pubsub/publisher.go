package pubsub

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Publisher interface {
	PublishBlock([]byte) error
	StartEventHandler()
	Close()
}

type KafkaProducer struct {
	*kafka.Producer
}

func NewPublisher(url string) (Publisher, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": url})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{p}, nil
}

func (p *KafkaProducer) StartEventHandler() {
	for e := range p.Producer.Events() {
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

func (p *KafkaProducer) PublishBlock(blockData []byte) error {
	err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &blocksTopic, Partition: kafka.PartitionAny},
		Value:          blockData,
	}, nil)
	return err
}

func (k *KafkaProducer) Close() {
	i := k.Producer.Flush(10000)
	for i > 0 {
		slog.Info("flushing kafka messages...", "remaining:", i)
		i = k.Producer.Flush(10000)
	}
	k.Producer.Close()
}
