package db

import (
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	InsertBlock(data.Block) error
	GetBlockByNumber(uint64) (*data.Block, error)
	GetFirstBlock() (*data.Block, error)
	GetBlocks() ([]*data.Block, error)
	Close() error
}

type GormDB struct {
	*gorm.DB
}

func NewConnection(url string) (DB, error) {
	db, err := gorm.Open(postgres.Open(url),
		&gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &GormDB{db}, nil
}

func (g *GormDB) Close() error {
	db, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}
	return db.Close()
}

func runMigrations(g *gorm.DB) error {
	err := g.AutoMigrate(&data.Block{})
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	err = g.AutoMigrate(&data.Transaction{})
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	return nil
}
