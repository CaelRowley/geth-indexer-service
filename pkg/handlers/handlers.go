package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
)

type Handlers struct {
	dbConn    db.DB
	ethClient *ethclient.Client
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func Init(dbConn db.DB, ethClient *ethclient.Client, r *chi.Mux) {
	handlers := Handlers{
		dbConn:    dbConn,
		ethClient: ethClient,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the API"))
	})

	r.Route("/block", func(r chi.Router) {
		handlers.AddBlockHandlers(r)
	})
}

func makeHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				setJSONResponse(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "interal server error",
				}
				setJSONResponse(w, http.StatusInternalServerError, errResp)
			}
			log.Println("HTTP API error:", err.Error(), "path:", r.URL.Path)
			// TODO: setup slog.Error
		}
	}
}

func setJSONResponse(w http.ResponseWriter, code int, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	w.Write(data)

	return nil
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJson(msg ...error) APIError {
	var err error
	if len(msg) > 0 {
		err = fmt.Errorf("invalid JSON request data %v", msg[0])
	} else {
		err = fmt.Errorf("invalid JSON request data")
	}
	return NewAPIError(http.StatusBadRequest, err)
}

func InvalidURLParam(msg ...error) APIError {
	var err error
	if len(msg) > 0 {
		err = fmt.Errorf("invalid URLParam %w", msg[0])
	} else {
		err = fmt.Errorf("invalid URLParam")
	}
	return NewAPIError(http.StatusBadRequest, err)
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}
