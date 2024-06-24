package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func StartSyncer(client *ethclient.Client, dbConn db.DB) {
	var nextBlockNumber uint64

	firstBlock, err := db.GetFirstBlock(dbConn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			block, err := client.BlockByNumber(context.Background(), nil)
			if err != nil {
				log.Fatalf("failed to retrieve the latest block from eth client: %v", err)
			}
			nextBlockNumber = block.NumberU64()
		} else {
			log.Fatalf("failed to retrieve the latest block from db: %v", err)
		}
	} else {
		nextBlockNumber = firstBlock.Number - 1
	}

	fmt.Println("Syncing blocks...")

	for nextBlockNumber > 0 {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(nextBlockNumber))
		if err != nil {
			log.Fatalf("failed to retrieve block number %d: %w", nextBlockNumber, err)
		}
		insertBlock(dbConn, block)
		nextBlockNumber -= 1
	}
}
