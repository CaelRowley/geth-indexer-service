package eth

import (
	"context"
	"fmt"
	"log"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
)

func StartListener(client *ethclient.Client, dbConn *pgx.Conn) {
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
			insertBlock(client, dbConn, header.Hash())
		}
	}
}

func insertBlock(client *ethclient.Client, dbConn *pgx.Conn, blockHash common.Hash) {
	block, err := client.BlockByHash(context.Background(), blockHash)
	if err != nil {
		log.Fatalf("failed to get block by hash: %v", err)
	}

	newBlock := data.Block{
		Hash:        block.Hash().Hex(),
		Number:      block.Number().Uint64(),
		GasLimit:    block.GasLimit(),
		GasUsed:     block.GasUsed(),
		Difficulty:  block.Difficulty().String(),
		Time:        block.Time(),
		ParentHash:  block.ParentHash().Hex(),
		Nonce:       hexutil.EncodeUint64(block.Nonce()),
		Miner:       block.Coinbase().Hex(),
		Size:        block.Size(),
		RootHash:    block.Root().Hex(),
		UncleHash:   block.UncleHash().Hex(),
		TxHash:      block.TxHash().Hex(),
		ReceiptHash: block.ReceiptHash().Hex(),
		ExtraData:   block.Extra(),
	}

	err = db.InsertBlock(context.Background(), dbConn, newBlock)
	if err != nil {
		log.Fatalf("failed to insert block: %v", err)
	}
}
