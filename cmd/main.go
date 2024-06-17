package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	createDB()
	connectToNode()
}

func createDB() {
	dbURL := "postgresql://admin:admin@localhost:5432"
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = createTables(conn)
	if err != nil {
		log.Fatal(err)
	}
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

func connectToNode() {
	nodeURL := os.Getenv("NODE_URL")

	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}

	ch := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Fatalf("Failed to subscribe to head: %v", err)
	}

	fmt.Println("Listening for new blocks...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-ch:
			fmt.Printf("New block: #%d (hash: %s)\n", header.Number.Uint64(), header.Hash().Hex())
			printBlockDetails(client, header.Hash())
			fmt.Println()
		}
	}
}

func printBlockDetails(client *ethclient.Client, blockHash common.Hash) {
	block, err := client.BlockByHash(context.Background(), blockHash)
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	fmt.Printf("Block details:\n")
	fmt.Printf("  Number: %d\n", block.Number().Uint64())
	fmt.Printf("  Hash: %s\n", block.Hash().Hex())
	fmt.Printf("  TxHash: %s\n", block.TxHash())
	fmt.Printf("  Transactions: %d\n", len(block.Transactions()))
	fmt.Printf("  GasUsed: %d\n", block.GasUsed())
}
