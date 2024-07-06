package eth

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
			block, err := c.BlockByNumber(context.Background(), new(big.Int).SetUint64(nextBlockNumber))
			if err != nil {
				slog.Error("syncer failed to get block", "number", nextBlockNumber, "err", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err = c.publishBlock(block); err != nil {
				if err.(kafka.Error).Code() == kafka.ErrQueueFull {
					time.Sleep(time.Second)
					continue
				}
				slog.Error("syncer failed to publish block", "number", nextBlockNumber, "err", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			nextBlockNumber -= 1
		}
	}
	return nil
}
