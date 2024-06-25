package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/eth"
	"github.com/CaelRowley/geth-indexer-service/pkg/router"
	"github.com/ethereum/go-ethereum/ethclient"
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

	r := router.NewRouter(dbConn)

	s := &Server{
		router:    r,
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

	errCh := make(chan error, 1)

	go func() {
		log.Println("Server be jammin' on port:", s.port)
		err := httpServer.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		close(errCh)
	}()

	if s.sync {
		log.Println("Syncing blocks on node with db...")
		go func() {
			err := eth.StartListener(s.ethClient, s.dbConn)
			if err != nil {
				errCh <- fmt.Errorf("failed to start listener: %w", err)
			}
		}()

		go func() {
			err := eth.StartSyncer(s.ethClient, s.dbConn)
			if err != nil {
				errCh <- fmt.Errorf("failed to start syncer: %w", err)
			}
		}()
	}

	defer func() {
		db, err := s.dbConn.DB()
		if err != nil {
			log.Println("failed to close db:", err)
			return
		}
		if err := db.Close(); err != nil {
			log.Println("failed to close db:", err)
			return
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
