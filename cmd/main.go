package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/CaelRowley/geth-indexer-service/cmd/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file:", err)
	}

	server, err := server.New()
	if err != nil {
		log.Fatal("failed to create server:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = server.Start(ctx)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
