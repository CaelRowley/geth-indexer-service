package handlers

import (
	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Handlers struct {
	dbConn    db.DB
	ethClient *ethclient.Client
}

func Init(dbConn db.DB, ethClient *ethclient.Client) *Handlers {
	return &Handlers{
		dbConn:    dbConn,
		ethClient: ethClient,
	}
}
