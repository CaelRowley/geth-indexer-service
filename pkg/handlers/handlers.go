package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	dbConn *pgx.Conn
}

func NewHandler(dbConn *pgx.Conn) *Handler {
	return &Handler{
		dbConn: dbConn,
	}
}

func (h *Handler) GetBlock(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	block, err := db.GetBlockByNumber(context.Background(), h.dbConn, number)

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("no block found with id:", number)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(block)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
