package router

import (
	"net/http"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/handlers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(dbConn db.DB, ethClient *ethclient.Client) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	handlers.Init(dbConn, ethClient, r)

	return r
}
