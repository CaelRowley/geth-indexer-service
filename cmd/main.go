package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/CaelRowley/geth-indexer-service/cmd/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "err", err)
		os.Exit(1)
	}

	var serverCfg server.ServerConfig
	flag.StringVar(&serverCfg.Port, "port", "8080", "Port where the service will run")
	flag.BoolVar(&serverCfg.Sync, "sync", false, "Sync blocks on node with db")
	flag.Parse()

	slog.Info("flags set", "Port", serverCfg.Port, "Sync", serverCfg.Sync)

	s, err := server.New(serverCfg)
	if err != nil {
		slog.Error("failed to create server", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := s.Start(ctx); err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
