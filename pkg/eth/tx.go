package eth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

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

func (c EthClient) publishTxs(ctx context.Context, txs types.Transactions, blockHash common.Hash) error {
	receipts, err := c.batchTransactionReceipts(ctx, txs)
	if err != nil {
		return err
	}
	if len(receipts) != len(txs) {
		return fmt.Errorf("len of receipts: %d doesnt match len of txs: %d", len(receipts), len(txs))
	}

	senders, err := c.batchTransactionSenders(ctx, txs, blockHash, receipts)
	if err != nil {
		return err
	}
	if len(senders) != len(txs) {
		return fmt.Errorf("len of senders: %d doesnt match len of txs: %d", len(senders), len(txs))
	}

	for i, tx := range txs {
		if err := c.publishTx(tx, senders[i], receipts[i]); err != nil {
			return err
		}
	}

	return nil
}

func (c *EthClient) batchTransactionReceipts(ctx context.Context, txs []*types.Transaction) ([]*types.Receipt, error) {
	var reqs []rpc.BatchElem
	for _, tx := range txs {
		reqs = append(reqs, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx.Hash()},
			Result: new(types.Receipt),
		})
	}

	if err := c.Client.Client().BatchCallContext(ctx, reqs); err != nil {
		return nil, err
	}

	receipts := make([]*types.Receipt, len(reqs))
	for i, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
		receipt, ok := req.Result.(*types.Receipt)
		if !ok {
			return nil, fmt.Errorf("unexpected type for tx receipt: %T", req.Result)
		}
		receipts[i] = receipt
	}

	return receipts, nil
}

func (c *EthClient) batchTransactionSenders(ctx context.Context, txs []*types.Transaction, blockHash common.Hash, receipts []*types.Receipt) ([]common.Address, error) {
	senders := make([]common.Address, len(txs))
	for i, tx := range txs {
		sender, err := c.TransactionSender(ctx, tx, blockHash, receipts[i].TransactionIndex)
		if err != nil {
			return nil, err
		}
		senders[i] = sender
	}

	return senders, nil
}
