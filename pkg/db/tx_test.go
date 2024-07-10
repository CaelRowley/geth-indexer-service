package db

import (
	"regexp"
	"testing"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
		Status:    0,
		BlockHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
	},
	{
		Hash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		From: "0x0000000000000000000000000000000000000003",
		// To:        ,
		Contract:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		Value:     200,
		Data:      []byte("some other data"),
		Gas:       500000,
		GasPrice:  1000000,
		Cost:      1000000000,
		Nonce:     0,
		Status:    0,
		BlockHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
	},
}

func TestInsertTx(t *testing.T) {
	s := newSuite(t)

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "transactions" ("hash","from","to","contract","value","data","gas","gas_price","cost","nonce","status","block_hash") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`)).
		WithArgs(
			mockTxs[0].Hash, mockTxs[0].From, mockTxs[0].To, mockTxs[0].Contract, mockTxs[0].Value, mockTxs[0].Data,
			mockTxs[0].Gas, mockTxs[0].GasPrice, mockTxs[0].Cost, mockTxs[0].Nonce, mockTxs[0].Status, mockTxs[0].BlockHash,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.sqlMock.ExpectCommit()
	err := s.dbMock.InsertTx(mockTxs[0])
	assert.NoError(t, err)
	assert.NoError(t, s.sqlMock.ExpectationsWereMet())
}

func TestGetTxByHash(t *testing.T) {
	s := newSuite(t)

	s.sqlMock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE hash = $1 ORDER BY "transactions"."hash" LIMIT $2`)).
		WithArgs(mockTxs[0].Hash, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"hash", "from", "to", "contract", "value", "data", "gas", "gas_price", "cost", "nonce", "status", "block_hash",
		}).AddRow(
			mockTxs[0].Hash, mockTxs[0].From, mockTxs[0].To, mockTxs[0].Contract, mockTxs[0].Value, mockTxs[0].Data,
			mockTxs[0].Gas, mockTxs[0].GasPrice, mockTxs[0].Cost, mockTxs[0].Nonce, mockTxs[0].Status, mockTxs[0].BlockHash,
		))

	retrievedBlock, err := s.dbMock.GetTxByHash(mockTxs[0].Hash)
	assert.NoError(t, err)
	assert.Equal(t, &mockTxs[0], retrievedBlock)
	assert.NoError(t, s.sqlMock.ExpectationsWereMet())
}

func TestGetTxs(t *testing.T) {
	s := newSuite(t)

	s.sqlMock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions"`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"hash", "from", "to", "contract", "value", "data", "gas", "gas_price", "cost", "nonce", "status", "block_hash",
		}).AddRow(
			mockTxs[0].Hash, mockTxs[0].From, mockTxs[0].To, mockTxs[0].Contract, mockTxs[0].Value, mockTxs[0].Data,
			mockTxs[0].Gas, mockTxs[0].GasPrice, mockTxs[0].Cost, mockTxs[0].Nonce, mockTxs[0].Status, mockTxs[0].BlockHash,
		).AddRow(
			mockTxs[1].Hash, mockTxs[1].From, mockTxs[1].To, mockTxs[1].Contract, mockTxs[1].Value, mockTxs[1].Data,
			mockTxs[1].Gas, mockTxs[1].GasPrice, mockTxs[1].Cost, mockTxs[1].Nonce, mockTxs[1].Status, mockTxs[1].BlockHash,
		))

	retrievedBlocks, err := s.dbMock.GetTxs()
	assert.NoError(t, err)
	assert.Len(t, retrievedBlocks, len(mockTxs))
	for i, retrievedBlock := range retrievedBlocks {
		assert.Equal(t, &mockTxs[i], retrievedBlock)
	}
	assert.NoError(t, s.sqlMock.ExpectationsWereMet())
}
