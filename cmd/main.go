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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	var serverCfg server.ServerConfig
	flag.BoolVar(&serverCfg.Sync, "sync", false, "Sync blocks on node with db")
	flag.StringVar(&serverCfg.Port, "port", "8080", "Port where the service will run")
	flag.Parse()

	s, err := server.New(serverCfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Start(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
