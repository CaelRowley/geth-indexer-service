package eth

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
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
			_, err := c.handleBlock(ctx, new(big.Int).SetUint64(nextBlockNumber))
			if err != nil {
				slog.Error("syncer failed to publish block", "err", err)
				time.Sleep(time.Millisecond * 100)
				continue
			}
			// for _, tx := range block.Transactions() {
			// 	if err := c.handleTx(ctx, tx, block.Hash()); err != nil {
			// 		slog.Error("syncer failed to publish tx", "err", err)
			// 		continue
			// 	}
			// }

			nextBlockNumber -= 1
		}
	}
	return nil
}

func (c EthClient) handleBlock(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := c.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if err = c.publishBlock(block); err != nil {
		return nil, err
	}
	return block, nil
}
