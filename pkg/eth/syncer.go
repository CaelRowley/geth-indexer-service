package eth

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"gorm.io/gorm"
)

func (c EthClient) StartSyncer(ctx context.Context, dbConn db.DB) error {
	var nextBlockNumber uint64
	firstBlock, err := dbConn.GetFirstBlock()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			block, err := c.BlockByNumber(context.Background(), nil)
			if err != nil {
				return fmt.Errorf("failed to retrieve the latest block from eth client: %w", err)
			}
			nextBlockNumber = block.NumberU64()
		} else {
			return fmt.Errorf("failed to retrieve the latest block from db: %w", err)
		}
	} else {
		nextBlockNumber = firstBlock.Number - 1
	}

	for nextBlockNumber > 0 {
		select {
		case <-ctx.Done():
			slog.Info("eth syncer stopped")
			return nil
		default:
			if err := c.handleBlock(ctx, nextBlockNumber); err != nil {
				slog.Error("syncer failed to publish block", "err", err)
				time.Sleep(time.Millisecond * 100)
				continue
			}
			nextBlockNumber -= 1
		}
	}
	return nil
}

func (c EthClient) handleBlock(ctx context.Context, number uint64) error {
	bigNumber := new(big.Int).SetUint64(number)
	block, err := c.BlockByNumber(ctx, bigNumber)
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
