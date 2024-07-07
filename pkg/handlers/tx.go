package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Handlers) GetTx(w http.ResponseWriter, r *http.Request) error {
	hash := chi.URLParam(r, "hash")
	tx, err := h.dbConn.GetTxByHash(hash)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}
	return setJSONResponse(w, http.StatusOK, tx)
}

func (h *Handlers) GetTxs(w http.ResponseWriter, r *http.Request) error {
	txs, err := h.dbConn.GetTxs()
	if err != nil {
		return fmt.Errorf("failed to get blocks: %w", err)
	}
	return setJSONResponse(w, http.StatusOK, txs)
}
