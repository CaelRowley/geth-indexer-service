package eth

import (
	"context"
	"encoding/json"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c EthClient) handleTx(ctx context.Context, tx *types.Transaction, blockHash common.Hash) error {
	receipt, err := c.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return err
	}
	sender, err := c.TransactionSender(ctx, tx, blockHash, receipt.TransactionIndex)
	if err != nil {
		return err
	}
	if err = c.publishTx(tx, sender, receipt); err != nil {
		return err
	}
	return nil
}

func (c EthClient) publishTx(tx *types.Transaction, sender common.Address, receipt *types.Receipt) error {
	newTx := data.Transaction{
		Hash:      tx.Hash().Hex(),
		From:      sender.Hex(),
		Contract:  receipt.ContractAddress.Hex(),
		Value:     tx.Value().Uint64(),
		Data:      tx.Data(),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice().Uint64(),
		Cost:      tx.Cost().Uint64(),
		Nonce:     tx.Nonce(),
		Status:    receipt.Status,
		BlockHash: receipt.BlockHash.Hex(),
	}
	if tx.To() != nil {
		newTx.To = tx.To().Hex()
	}

	txData, err := json.Marshal(newTx)
	if err != nil {
		return err
	}
	if err := c.PubSub.GetPublisher().PublishTx(txData); err != nil {
		return err
	}
	return nil
}
