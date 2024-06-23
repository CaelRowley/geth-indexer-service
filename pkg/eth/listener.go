package eth

import (
	"context"
	"fmt"
	"log"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener(client *ethclient.Client, dbConn db.DB) {
	ch := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Fatalf("failed to subscribe to head: %v", err)
	}

	fmt.Println("Listening for new blocks...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case header := <-ch:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatalf("failed to get block by hash: %v", err)
			}
			insertBlock(dbConn, block)
		}
	}
}
