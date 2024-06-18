package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	db *pgx.Conn
}

type Block struct {
	ID     int64  `json:"id"`
	Number int64  `json:"number"`
	Hash   string `json:"hash"`
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetBlock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := `SELECT id, number, hash FROM blocks WHERE id = $1`

	var block Block

	err := h.db.QueryRow(context.Background(), query, id).Scan(&block.ID, &block.Number, &block.Hash)
	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("no block found with id:", id)
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
