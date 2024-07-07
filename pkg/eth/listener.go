package eth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c EthClient) StartListener(ctx context.Context, dbConn db.DB) error {
	headerCh := make(chan *types.Header)
	defer close(headerCh)
	sub, err := c.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe to head: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err) // TODO: add retry logic
		case header := <-headerCh:
			if err := c.handleHeader(ctx, header); err != nil {
				slog.Error("failed to process header", "err", err)
			}
		case <-ctx.Done():
			slog.Info("eth listener stopped")
			return nil
		}
	}
}

func (c EthClient) handleHeader(ctx context.Context, header *types.Header) error {
	block, err := c.BlockByHash(ctx, header.Hash())
	if err != nil {
		return err
	}
	if err = c.publishBlock(block); err != nil {
		return err
	}
	if err := c.publishTxs(ctx, block.Transactions(), block.Hash()); err != nil {
		return err
	}
	return nil
}
