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
			slog.Error("syncer failed to get block", "number", nextBlockNumber, "err", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err = c.publishBlock(block)
		if err != nil {
			slog.Error("syncer failed to publish block", "number", nextBlockNumber, "err", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// newBlock := data.Block{
		// 	Hash:        block.Hash().Hex(),
		// 	Number:      block.Number().Uint64(),
		// 	GasLimit:    block.GasLimit(),
		// 	GasUsed:     block.GasUsed(),
		// 	Difficulty:  block.Difficulty().String(),
		// 	Time:        block.Time(),
		// 	ParentHash:  block.ParentHash().Hex(),
		// 	Nonce:       hexutil.EncodeUint64(block.Nonce()),
		// 	Miner:       block.Coinbase().Hex(),
		// 	Size:        block.Size(),
		// 	RootHash:    block.Root().Hex(),
		// 	UncleHash:   block.UncleHash().Hex(),
		// 	TxHash:      block.TxHash().Hex(),
		// 	ReceiptHash: block.ReceiptHash().Hex(),
		// 	ExtraData:   block.Extra(),
		// }

		// blockData, err := json.Marshal(newBlock)
		// if err != nil {
		// 	slog.Error("failed to serialize block data", "err", err)
		// 	time.Sleep(100 * time.Millisecond)
		// 	continue
		// }

		// if err := c.PubSub.PublishBlock(blockData); err != nil {
		// 	slog.Error("failed to publish block", "err", err)
		// 	time.Sleep(100 * time.Millisecond)
		// 	continue
		// }

		nextBlockNumber -= 1
	}

	return nil
}
