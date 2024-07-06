package eth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c EthClient) StartListener(ctx context.Context, dbConn db.DB) error {
	headerCh := make(chan *types.Header)
	defer close(headerCh)
	sub, err := c.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		return fmt.Errorf("failed to subscribe to head: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err) // TODO: add retry logic
		case header := <-headerCh:
			header.Hash()
			if err := c.processHeader(ctx, header); err != nil {
				slog.Error("failed to process header", "err", err)
			}

		case <-ctx.Done():
			slog.Info("listener stopped")
			return nil
		}
	}
}

func (c EthClient) processHeader(ctx context.Context, header *types.Header) error {
	block, err := c.BlockByHash(ctx, header.Hash())
	if err != nil {
		return fmt.Errorf("failed to get block by hash %s: %w", header.Hash().Hex(), err)
	}

	return c.publishBlock(block)

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
	// 	return fmt.Errorf("failed to serialize block data: %w", err)
	// }

	// if err := client.PubSub.PublishBlock(blockData); err != nil {
	// 	return fmt.Errorf("failed to publish block: %w", err)
	// }

	// return nil
}
