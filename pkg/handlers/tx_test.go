package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

var mockTxs = []data.Transaction{
	{
		Hash:      "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		From:      "0x0000000000000000000000000000000000000001",
		To:        "0x0000000000000000000000000000000000000002",
		Contract:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		Value:     200,
		Data:      []byte("some data"),
		Gas:       500000,
		GasPrice:  1000000,
		Cost:      1000000000,
		Nonce:     0,
		Status:    1,
		BlockHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
	},
	{
		Hash:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		From:      "0x0000000000000000000000000000000000000003",
		To:        "",
		Contract:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Value:     500,
		Data:      []byte("some other data"),
		Gas:       600000,
		GasPrice:  2000000,
		Cost:      2000000000,
		Nonce:     1,
		Status:    1,
		BlockHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	},
}

func TestGetTx(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetTxByHash", mockTxs[0].Hash).Return(mockTxs[0], nil)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-tx/{hash}", makeHandler(handlers.GetTx))

	req, err := http.NewRequest("GET", "/get-tx/"+mockTxs[0].Hash, nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var tx data.Transaction
	err = json.NewDecoder(recorder.Body).Decode(&tx)
	assert.NoError(t, err)
	assert.Equal(t, mockTxs[0], tx)

	mockDB.AssertExpectations(t)
}

func TestGetTxs(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetTxs").Return([]*data.Transaction{&mockTxs[0], &mockTxs[1]}, nil)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-txs", makeHandler(handlers.GetTxs))

	req, err := http.NewRequest("GET", "/get-txs", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var txs []data.Transaction
	err = json.NewDecoder(recorder.Body).Decode(&txs)
	assert.NoError(t, err)
	assert.Len(t, txs, len(mockTxs))
	assert.Equal(t, mockTxs, txs)

	mockDB.AssertExpectations(t)
}

func TestGetTxDBError(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetTxByHash", "0xnonexistent").Return(nil, errors.New("tx not found"))

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-tx/{hash}", makeHandler(handlers.GetTx))

	req, err := http.NewRequest("GET", "/get-tx/0xnonexistent", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockDB.AssertExpectations(t)
}

func TestGetTxsDBError(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetTxs").Return([]*data.Transaction(nil), errors.New("database error"))

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-txs", makeHandler(handlers.GetTxs))

	req, err := http.NewRequest("GET", "/get-txs", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockDB.AssertExpectations(t)
}

func TestGetBlockInvalidNumberParam(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-block/{number}", makeHandler(handlers.GetBlock))

	req, err := http.NewRequest("GET", "/get-block/invalid", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetBlockDBError(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetBlockByNumber", uint64(999)).Return(nil, errors.New("block not found"))

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-block/{number}", makeHandler(handlers.GetBlock))

	req, err := http.NewRequest("GET", "/get-block/999", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockDB.AssertExpectations(t)
}

func TestGetBlocksDBError(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetBlocks").Return([]*data.Block(nil), errors.New("database error"))

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-blocks", makeHandler(handlers.GetBlocks))

	req, err := http.NewRequest("GET", "/get-blocks", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockDB.AssertExpectations(t)
}
