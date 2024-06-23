package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB = *pgx.Conn

func NewConnection(url string) (DB, error) {
	dbConn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// TODO: setup migrations for table creation
	err = createTables(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return dbConn, nil
}

func createTables(dbConn DB) error {
	dropTable := `DROP TABLE IF EXISTS blocks`
	_, err := dbConn.Exec(context.Background(), dropTable)
	if err != nil {
		return err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS blocks (
				hash TEXT NOT NULL,
				number NUMERIC NOT NULL,
				gas_limit    NUMERIC NOT NULL,
				gas_used     NUMERIC NOT NULL,
				difficulty  TEXT NOT NULL,
				time        NUMERIC NOT NULL,
				parent_hash  TEXT NOT NULL,
				nonce       TEXT NOT NULL,
				miner       TEXT NOT NULL,
				size        NUMERIC NOT NULL,
				root_hash    TEXT NOT NULL,
				uncle_hash   TEXT NOT NULL,
				tx_hash      TEXT NOT NULL,
				receipt_hash TEXT NOT NULL,
				extra_data   BYTEA NOT NULL
		)
	`

	_, err = dbConn.Exec(context.Background(), createTableQuery)
	if err != nil {
		return err
	}

	return nil
}
