//go:build !kafka

package pubsub

import (
	"context"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
)

type Subscriber interface {
	StartPoll(context.Context) error
	Close() error
}

type StubSubscriber struct{}

func NewSubscriber(url string, dbConn db.DB) (Subscriber, error) {
	return &StubSubscriber{}, nil
}

func (s *StubSubscriber) StartPoll(ctx context.Context) error {
	return nil
}

func (s *StubSubscriber) Close() error {
	return nil
}