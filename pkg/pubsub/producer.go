package pubsub

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer(url string) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": url})
}

func (k *KafkaPubSub) PublishBlock(blockData []byte) error {
	fmt.Println("publishing block")
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &blocksTopic, Partition: kafka.PartitionAny},
		Value:          blockData,
	}, nil)

	return err
}