package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/eth"
	"github.com/CaelRowley/geth-indexer-service/pkg/handlers"
	"github.com/CaelRowley/geth-indexer-service/pkg/pubsub"
	"github.com/CaelRowley/geth-indexer-service/pkg/router"
	"golang.org/x/exp/slog"
)

type ServerConfig struct {
	Sync bool
	Port string
}

type Server struct {
	router    http.Handler
	dbConn    db.DB
	ethClient eth.Client
	pubsub    pubsub.PubSub
	sync      bool
	port      string
}

func New(cfg ServerConfig) (*Server, error) {
	dbConn, err := db.NewConnection(os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	pubsubClient, err := pubsub.NewPubSub(os.Getenv("MSG_BROKER_URL"), dbConn)
	if err != nil {
		return nil, err
	}
	ethClient, err := eth.NewClient(os.Getenv("NODE_URL"), pubsubClient)
	if err != nil {
		return nil, err
	}
	router := router.NewRouter()
	handlers.Init(dbConn, router)

	s := &Server{
		router:    router,
		dbConn:    dbConn,
		ethClient: ethClient,
		pubsub:    pubsubClient,
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
	var wg sync.WaitGroup

	go func() {
		slog.Info("server be jammin' on", "port", s.port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http server failed %w", err)
		}
	}()

	if s.sync {
		slog.Info("starting sync with eth node...")
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.ethClient.StartListener(ctx, s.dbConn); err != nil {
				errCh <- fmt.Errorf("listener failed: %w", err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.ethClient.StartSyncer(ctx, s.dbConn); err != nil {
				errCh <- fmt.Errorf("syncer failed: %w", err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.pubsub.GetSubscriber().StartPoll(ctx); err != nil {
				errCh <- fmt.Errorf("consumer failed: %w", err)
			}
		}()
		go s.pubsub.GetPublisher().StartEventHandler()
	}

	defer func() {
		if err := s.pubsub.Close(); err != nil {
			slog.Error("failed to close pubsub", "err", err)
		}
		s.ethClient.Close()
		if err := s.dbConn.Close(); err != nil {
			slog.Error("failed to close db", "err", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		wg.Wait()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return httpServer.Shutdown(timeoutCtx)
	}
}
