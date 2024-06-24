package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func StartSyncer(client *ethclient.Client, dbConn db.DB) error {
	var nextBlockNumber uint64

	firstBlock, err := db.GetFirstBlock(dbConn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			block, err := client.BlockByNumber(context.Background(), nil)
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

	fmt.Println("Syncing blocks...")

	for nextBlockNumber > 0 {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(nextBlockNumber))
		if err != nil {
			log.Printf("failed to get block by number %d: %v\n", nextBlockNumber, err)
			time.Sleep(500 * time.Millisecond)
		}
		err = insertBlock(dbConn, block)
		if err != nil {
			log.Printf("failed to insert block: number %d: %v\n", block.NumberU64(), err)
			time.Sleep(500 * time.Millisecond)
		}
		nextBlockNumber -= 1
	}

	return nil
}
