package eth

import (
	"context"
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/pubsub"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client interface {
	StartSyncer(db.DB) error
	StartListener(context.Context, db.DB) error
}

type EthClient struct {
	*ethclient.Client
	pubsub.PubSub
}

func NewClient(url string, pubsub pubsub.PubSub) (Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth client: %w", err)
	}

	return &EthClient{client, pubsub}, nil
}
