package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewConnection(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// TODO: setup migrations for table creation
	err = createTables(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return conn, nil
}

func createTables(conn *pgx.Conn) error {
	dropTable := `DROP TABLE IF EXISTS blocks`
	_, err := conn.Exec(context.Background(), dropTable)
	if err != nil {
		return err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS blocks (
				id SERIAL PRIMARY KEY,
				Number BIGINT NOT NULL,
				Hash TEXT NOT NULL
		)
	`

	_, err = conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		return err
	}

	return nil
}
