package db

import (
	"context"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/jackc/pgx/v5"
)

func InsertBlock(ctx context.Context, dbConn *pgx.Conn, block data.Block) error {
	query := `
		INSERT INTO blocks (
			hash, number, gas_limit, gas_used, difficulty, time, parent_hash, nonce, miner, size, root_hash, uncle_hash, tx_hash, receipt_hash, extra_data
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err := dbConn.Exec(ctx, query,
		block.Hash,
		block.Number,
		block.GasLimit,
		block.GasUsed,
		block.Difficulty,
		block.Time,
		block.ParentHash,
		block.Nonce,
		block.Miner,
		block.Size,
		block.RootHash,
		block.UncleHash,
		block.TxHash,
		block.ReceiptHash,
		block.ExtraData,
	)

	return err
}
