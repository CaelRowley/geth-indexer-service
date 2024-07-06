package pubsub

import (
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
)

var (
	blocksTopic = "blocks"
)

type PubSub interface {
	GetPublisher() Publisher
	GetSubscriber() Subscriber
	Close() error
}

type Broker struct {
	Publisher
	Subscriber
}

func NewPubSub(url string, dbConn db.DB) (PubSub, error) {
	p, err := NewPublisher(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}
	s, err := NewSubscriber(url, dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	return &Broker{p, s}, nil
}

func (b *Broker) GetPublisher() Publisher {
	return b.Publisher
}

func (b *Broker) GetSubscriber() Subscriber {
	return b.Subscriber
}

func (b *Broker) Close() error {
	b.Publisher.Close()
	return b.Subscriber.Close()
}
