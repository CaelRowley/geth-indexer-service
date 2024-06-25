package router

import (
	"net/http"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/CaelRowley/geth-indexer-service/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter(dbConn db.DB) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := handlers.NewHandler(dbConn)

	router.Route("/eth", func(r chi.Router) {
		loadEthRoutes(r, *h)
	})

	return router
}

func loadEthRoutes(router chi.Router, h handlers.Handler) {
	router.Get("/get-block/{number}", h.GetBlock)
	router.Get("/get-blocks", h.GetBlocks)
}
