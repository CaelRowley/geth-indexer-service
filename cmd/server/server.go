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

type Server struct {
	router    http.Handler
	dbConn    db.DB
	ethClient *ethclient.Client
}

var port = "8080"

func New() *Server {
	dbConn, err := db.NewConnection(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ethClient, err := eth.NewClient(os.Getenv("NODE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	router := router.NewRouter(dbConn)

	server := &Server{
		router:    router,
		dbConn:    dbConn,
		ethClient: ethClient,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	defer func() {
		if err := s.dbConn.Close(context.Background()); err != nil {
			fmt.Println("failed to close db", err)
		}
	}()

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	fmt.Println("Server be jammin on port:", port)

	go eth.StartSyncer(s.ethClient, s.dbConn)
	go eth.StartListener(s.ethClient, s.dbConn)

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
