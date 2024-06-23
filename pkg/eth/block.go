package eth

import (
	"context"
	"log"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func insertBlock(dbConn db.DB, block *types.Block) {
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

	err := db.InsertBlock(context.Background(), dbConn, newBlock)
	if err != nil {
		log.Fatalf("failed to insert block: %v", err)
	}
}
