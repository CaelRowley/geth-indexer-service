package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/go-chi/chi"
)

func (h *Handlers) AddBlockHandlers(r chi.Router) {
	r.Get("/get-block/{number}", makeHandler(h.GetBlock))
	r.Get("/get-blocks", makeHandler(h.GetBlocks))
}

func (h *Handlers) GetBlock(w http.ResponseWriter, r *http.Request) error {
	number, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		return InvalidURLParam(fmt.Errorf("number: %w", err))
	}

	block, err := db.GetBlockByNumber(h.dbConn, number)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	return setJSONResponse(w, http.StatusOK, block)
}

func (h *Handlers) GetBlocks(w http.ResponseWriter, r *http.Request) error {
	blocks, err := db.GetBlocks(h.dbConn)
	if err != nil {
		return fmt.Errorf("failed to get blocks: %w", err)
	}

	return setJSONResponse(w, http.StatusOK, blocks)
}
