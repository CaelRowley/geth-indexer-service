package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/eth"
	"github.com/CaelRowley/geth-indexer-service/pkg/router"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/exp/slog"
)

type ServerConfig struct {
	Sync bool
	Port string
}

type Server struct {
	router    http.Handler
	dbConn    db.DB
	ethClient *ethclient.Client
	sync      bool
	port      string
}

func New(cfg ServerConfig) (*Server, error) {
	dbConn, err := db.NewConnection(os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	ethClient, err := eth.NewClient(os.Getenv("NODE_URL"))
	if err != nil {
		return nil, err
	}

	router := router.NewRouter(dbConn, ethClient)

	s := &Server{
		router:    router,
		dbConn:    dbConn,
		ethClient: ethClient,
		sync:      cfg.Sync,
		port:      cfg.Port,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	errCh := make(chan error)

	go func() {
		slog.Info("server be jammin' on", "port", s.port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http server failed %w", err)
		}
	}()

	if s.sync {
		slog.Info("syncing blocks on node with db...")
		go func() {
			if err := eth.StartListener(ctx, s.ethClient, s.dbConn); err != nil {
				errCh <- fmt.Errorf("listener failed: %w", err)
			}
		}()

		go func() {
			if err := eth.StartSyncer(s.ethClient, s.dbConn); err != nil {
				errCh <- fmt.Errorf("syncer failed: %w", err)
			}
		}()
	}

	defer func() {
		db, err := s.dbConn.DB()
		if err != nil {
			slog.Error("failed to get db connection", "err", err)
			return
		}
		if err := db.Close(); err != nil {
			slog.Error("failed to close db", "err", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return httpServer.Shutdown(timeoutCtx)
	}
}
