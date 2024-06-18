package router

import (
	"net/http"

	"github.com/CaelRowley/geth-indexer-service/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

func NewRouter(db *pgx.Conn) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := handlers.NewHandler(db)

	router.Route("/eth", func(r chi.Router) {
		loadEthRoutes(r, *h)
	})

	return router
}

func loadEthRoutes(router chi.Router, h handlers.Handler) {
	router.Get("/get-block/{id}", h.GetBlock)
}
