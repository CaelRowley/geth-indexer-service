package eth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener(ctx context.Context, client *ethclient.Client, dbConn db.DB) error {
	headerCh := make(chan *types.Header)
	defer close(headerCh)
	sub, err := client.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe to head: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err) // TODO: add retry logic
		case header := <-headerCh:
			if err := processHeader(ctx, client, dbConn, header); err != nil {
				slog.Error("failed to process header", "err", err)
			}
		case <-ctx.Done():
			slog.Info("listener stopped")
			return nil
		}
	}
}

func processHeader(ctx context.Context, client *ethclient.Client, dbConn db.DB, header *types.Header) error {
	block, err := client.BlockByHash(ctx, header.Hash())
	if err != nil {
		return fmt.Errorf("failed to get block by hash %s: %w", header.Hash().Hex(), err)
	}

	if err := insertBlock(dbConn, block); err != nil {
		return fmt.Errorf("failed to insert block (number %d): %w", block.NumberU64(), err)
	}

	return nil
}
