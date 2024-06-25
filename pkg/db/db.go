package db

import (
	"fmt"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB = *gorm.DB

func NewConnection(url string) (DB, error) {
	dbConn, err := gorm.Open(postgres.Open(url),
		&gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := runMigrations(dbConn); err != nil {
		return nil, err
	}

	return dbConn, nil
}

func runMigrations(dbConn DB) error {
	err := dbConn.AutoMigrate(&data.Block{})
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}
