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
		log.Fatal("Error loading .env file")
	}

	server := server.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
