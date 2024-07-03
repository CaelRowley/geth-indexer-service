package db

import (
	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

func (g *GormDB) InsertBlock(block data.Block) error {
	return g.Create(&block).Error
}

func (g *GormDB) GetBlockByNumber(number uint64) (*data.Block, error) {
	var block data.Block
	if err := g.First(&block, "number = ?", number).Error; err != nil {
		return nil, err
	}
	return &block, nil
}

func (g *GormDB) GetFirstBlock() (*data.Block, error) {
	var block data.Block
	if err := g.Order("number asc").First(&block).Error; err != nil {
		return nil, err
	}
	return &block, nil
}

func (g *GormDB) GetBlocks() ([]*data.Block, error) {
	var blocks []*data.Block
	if err := g.Find(&blocks).Error; err != nil {
		return nil, err
	}
	return blocks, nil
}
