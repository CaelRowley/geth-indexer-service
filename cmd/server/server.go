package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	dbURL := os.Getenv("DB_URL")
	nodeURL := os.Getenv("NODE_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	ethClient, err := eth.NewClient(nodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}

	router := newRouter()

	server := &Server{
		router:    router,
		db:        conn,
		ethClient: ethClient,
	}

	return server
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

	// TODO: setup migrations for table creation
	err = s.createTables()
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
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

func (s *Server) createTables() error {
	dropTable := `DROP TABLE IF EXISTS blocks`
	_, err := s.db.Exec(context.Background(), dropTable)
	if err != nil {
		return err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS blocks (
				id SERIAL PRIMARY KEY,
				Number BIGINT NOT NULL,
				Hash TEXT NOT NULL
		)
	`

	_, err = s.db.Exec(context.Background(), createTableQuery)
	if err != nil {
		return err
	}

	return nil
}
