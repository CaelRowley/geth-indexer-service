package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/go-chi/chi"
)

func (h *Handlers) GetBlock(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseUint(chi.URLParam(r, "number"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid number parameter", http.StatusBadRequest)
		log.Println("invalid number parameter:", err)
		return
	}

	block, err := db.GetBlockByNumber(h.dbConn, number)
	if err != nil {
		log.Println("failed to get block:", err)
		http.Error(w, "Failed to get block", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(block)
	if err != nil {
		log.Println("failed to marshal block:", err)
		http.Error(w, "Failed to serialize block data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handlers) GetBlocks(w http.ResponseWriter, r *http.Request) {
	blocks, err := db.GetBlocks(h.dbConn)
	if err != nil {
		log.Println("failed to get blocks:", err)
		http.Error(w, "Failed to get blocks", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(blocks)
	if err != nil {
		log.Println("failed to marshal blocks:", err)
		http.Error(w, "Failed to serialize blocks data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
