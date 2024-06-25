package main

import (
	"context"
	"flag"
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

	var serverConfig server.ServerConfig
	flag.BoolVar(&serverConfig.Sync, "sync", false, "Sync blocks with node")
	flag.StringVar(&serverConfig.Port, "port", "8080", "Port where the service will run")
	flag.Parse()

	server, err := server.New(serverConfig)
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
