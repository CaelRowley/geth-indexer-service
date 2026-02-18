package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) InsertBlock(block data.Block) error {
	args := m.Called(block)
	return args.Error(0)
}

func (m *MockDB) GetBlockByNumber(number uint64) (*data.Block, error) {
	args := m.Called(number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	block := args.Get(0).(data.Block)
	return &block, args.Error(1)
}

func (m *MockDB) GetFirstBlock() (*data.Block, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	block := args.Get(0).(data.Block)
	return &block, args.Error(1)
}

func (m *MockDB) GetBlocks() ([]*data.Block, error) {
	args := m.Called()
	return args.Get(0).([]*data.Block), args.Error(1)
}

func (m *MockDB) InsertTx(tx data.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockDB) GetTxByHash(hash string) (*data.Transaction, error) {
	args := m.Called(hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	block := args.Get(0).(data.Transaction)
	return &block, args.Error(1)
}

func (m *MockDB) GetTxs() ([]*data.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]*data.Transaction), args.Error(1)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

var mockBlocks = []data.Block{
	{
		Hash:        "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		Number:      1,
		GasLimit:    1000000,
		GasUsed:     500000,
		Difficulty:  1000000000,
		Time:        1625812800,
		ParentHash:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		Nonce:       0,
		Miner:       "0x0000000000000000000000000000000000000001",
		Size:        200,
		RootHash:    "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		UncleHash:   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		TxHash:      "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		ReceiptHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		ExtraData:   []byte("some extra data"),
	},
	{
		Hash:        "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		Number:      2,
		GasLimit:    2000000,
		GasUsed:     1000000,
		Difficulty:  2000000000,
		Time:        1625812900,
		ParentHash:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Nonce:       0,
		Miner:       "0x0000000000000000000000000000000000000002",
		Size:        400,
		RootHash:    "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		UncleHash:   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		TxHash:      "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		ReceiptHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		ExtraData:   []byte("some more extra data"),
	},
}

func TestGetBlockInvalidNumber(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-block/{number}", makeHandler(handlers.GetBlock))

	req, err := http.NewRequest("GET", "/get-block/notanumber", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetBlockDBError(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetBlockByNumber", uint64(999)).Return(nil, fmt.Errorf("record not found"))

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
	mockDB.On("GetBlocks").Return([]*data.Block(nil), fmt.Errorf("db connection error"))

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

func TestGetBlock(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetBlockByNumber", uint64(1)).Return(mockBlocks[0], nil)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-block/{number}", makeHandler(handlers.GetBlock))

	req, err := http.NewRequest("GET", "/get-block/1", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var block data.Block
	err = json.NewDecoder(recorder.Body).Decode(&block)
	assert.NoError(t, err)
	assert.Equal(t, mockBlocks[0], block)

	mockDB.AssertExpectations(t)
}

func TestGetBlocks(t *testing.T) {
	recorder := httptest.NewRecorder()
	mockDB := new(MockDB)
	mockDB.On("GetBlocks").Return([]*data.Block{&mockBlocks[0], &mockBlocks[1]}, nil)

	handlers := &Handlers{
		dbConn: mockDB,
	}

	r := chi.NewRouter()
	r.Get("/get-blocks", makeHandler(handlers.GetBlocks))

	req, err := http.NewRequest("GET", "/get-blocks", nil)
	assert.NoError(t, err)

	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var blocks []data.Block
	err = json.NewDecoder(recorder.Body).Decode(&blocks)
	assert.NoError(t, err)
	assert.Len(t, blocks, len(mockBlocks))
	assert.Equal(t, mockBlocks, blocks)

	mockDB.AssertExpectations(t)
}
