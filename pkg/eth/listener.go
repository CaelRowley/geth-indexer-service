package eth

import (
	"context"
	"fmt"
	"log"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener(client *ethclient.Client, dbConn db.DB) error {
	headerCh := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headerCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe to head: %w", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case header := <-headerCh:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("failed to get block by hash %s: %v\n", header.Hash().Hex(), err)
			}
			err = insertBlock(dbConn, block)
			if err != nil {
				log.Printf("failed to insert block: number %d: %v\n", block.NumberU64(), err)
			}
		}
	}
}
