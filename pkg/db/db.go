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
	dbConn.AutoMigrate(&data.Block{})

	return dbConn, nil
}
