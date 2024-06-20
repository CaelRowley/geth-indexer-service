package eth

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
)

func StartSyncer(client *ethclient.Client, dbConn *pgx.Conn) {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve the latest block: %v", err)
	}

	fmt.Println("Syncing blocks...")
	nextBlockNumber := block.Number()
	insertBlock(dbConn, block)

	for nextBlockNumber.Uint64() > 0 {
		nextBlockNumber = new(big.Int).Sub(nextBlockNumber, big.NewInt(1))
		block, err := client.BlockByNumber(context.Background(), nextBlockNumber)
		if err != nil {
			log.Fatalf("Failed to retrieve the latest block: %v", err)
		}
		insertBlock(dbConn, block)
	}
}
