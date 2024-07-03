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

func (c EthClient) StartSyncer(dbConn db.DB) error {
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
		block, err := c.BlockByNumber(context.Background(), new(big.Int).SetUint64(nextBlockNumber))
		if err != nil {
			slog.Error("failed to get block", "number", nextBlockNumber, "err", err)
			time.Sleep(500 * time.Millisecond)
		}
		err = insertBlock(dbConn, block)
		if err != nil {
			slog.Error("failed to insert block", "number", block.NumberU64(), "err", err)
			time.Sleep(500 * time.Millisecond)
		}
		nextBlockNumber -= 1
	}

	return nil
}
