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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	router    http.Handler
	db        *pgx.Conn
	ethClient *ethclient.Client
}

var port = "8080"

func New() *Server {
	router := newRouter()

	conn, err := db.NewConnection(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ethClient, err := eth.NewClient(os.Getenv("NODE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		router:    router,
		db:        conn,
		ethClient: ethClient,
	}
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	err := s.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	defer func() {
		if err := s.db.Close(context.Background()); err != nil {
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

	go eth.StartListener(s.ethClient)

	fmt.Println("Server be jammin on port:", port)

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}

func newRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
