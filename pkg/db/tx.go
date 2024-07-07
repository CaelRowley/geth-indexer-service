package db

import (
	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

func (g *GormDB) InsertTx(tx data.Transaction) error {
	return g.Create(&tx).Error
}

func (g *GormDB) GetTxByHash(hash string) (*data.Transaction, error) {
	var tx data.Transaction
	if err := g.First(&tx, "hash = ?", hash).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (g *GormDB) GetTxs() ([]*data.Transaction, error) {
	var txs []*data.Transaction
	if err := g.Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
