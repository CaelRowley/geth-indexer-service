package db

import (
	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

func GetBlockByNumber(dbConn DB, number uint64) (*data.Block, error) {
	var block data.Block
	if err := dbConn.First(&block, "number = ?", number).Error; err != nil {
		return nil, err
	}
	return &block, nil
}

func GetFirstBlock(dbConn DB) (*data.Block, error) {
	var block data.Block
	if err := dbConn.Order("number asc").First(&block).Error; err != nil {
		return nil, err
	}
	return &block, nil
}

func InsertBlock(dbConn DB, block data.Block) error {
	return dbConn.Create(&block).Error
}
