package eth

import (
	"encoding/json"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c EthClient) publishBlock(block *types.Block) error {
	newBlock := data.Block{
		Hash:        block.Hash().Hex(),
		Number:      block.Number().Uint64(),
		GasLimit:    block.GasLimit(),
		GasUsed:     block.GasUsed(),
		Difficulty:  block.Difficulty().Uint64(),
		Time:        block.Time(),
		ParentHash:  block.ParentHash().Hex(),
		Nonce:       block.Nonce(),
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
		return err
	}
	if err := c.PubSub.GetPublisher().PublishBlock(blockData); err != nil {
		return err
	}
	return nil
}
