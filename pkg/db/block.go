package db

import (
	"context"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

func GetBlockByNumber(ctx context.Context, dbConn DB, number uint64) (*data.Block, error) {
	var block data.Block
	if err := dbConn.First(&block, "number = ?", number).Error; err != nil {
		return nil, err
	}
	return &block, nil
}

func InsertBlock(ctx context.Context, dbConn DB, block data.Block) error {
	return dbConn.Create(&block).Error
}
