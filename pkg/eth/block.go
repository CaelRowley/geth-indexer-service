package eth

import (
	"encoding/json"
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c EthClient) publishBlock(block *types.Block) error {
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

	blockData, err := json.Marshal(newBlock)
	if err != nil {
		return fmt.Errorf("failed to serialize block data: %w", err)
	}

	if err := c.PubSub.PublishBlock(blockData); err != nil {
		return fmt.Errorf("failed to publish block: %w", err)
	}

	return nil
}
