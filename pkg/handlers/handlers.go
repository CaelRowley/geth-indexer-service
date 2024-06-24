package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/go-chi/chi"
)

type Handler struct {
	dbConn db.DB
}

func NewHandler(dbConn db.DB) *Handler {
	return &Handler{
		dbConn: dbConn,
	}
}

func (h *Handler) GetBlock(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	block, err := db.GetBlockByNumber(h.dbConn, number)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(block)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
