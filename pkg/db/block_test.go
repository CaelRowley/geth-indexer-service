package db

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/DATA-DOG/go-sqlmock"
)

var mockBlocks = []data.Block{
	{
		Hash:        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Number:      1,
		GasLimit:    1000000,
		GasUsed:     500000,
		Difficulty:  "1000000000",
		Time:        1625812800,
		ParentHash:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		Nonce:       "0x0000000000000001",
		Miner:       "0x0000000000000000000000000000000000000001",
		Size:        200,
		RootHash:    "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		UncleHash:   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		TxHash:      "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		ReceiptHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		ExtraData:   []byte("some extra data"),
	},
	{
		Hash:        "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		Number:      2,
		GasLimit:    2000000,
		GasUsed:     1000000,
		Difficulty:  "2000000000",
		Time:        1625812900,
		ParentHash:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Nonce:       "0x0000000000000002",
		Miner:       "0x0000000000000000000000000000000000000002",
		Size:        400,
		RootHash:    "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		UncleHash:   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		TxHash:      "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		ReceiptHash: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		ExtraData:   []byte("some more extra data"),
	},
}

func TestInsertBlock(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mDB.sqlmock.ExpectBegin()

	mDB.sqlmock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "blocks" ("hash","number","gas_limit","gas_used","difficulty","time","parent_hash","nonce","miner","size","root_hash","uncle_hash","tx_hash","receipt_hash","extra_data") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`)).
		WithArgs(
			mockBlocks[0].Hash,
			mockBlocks[0].Number,
			mockBlocks[0].GasLimit,
			mockBlocks[0].GasUsed,
			mockBlocks[0].Difficulty,
			mockBlocks[0].Time,
			mockBlocks[0].ParentHash,
			mockBlocks[0].Nonce,
			mockBlocks[0].Miner,
			mockBlocks[0].Size,
			mockBlocks[0].RootHash,
			mockBlocks[0].UncleHash,
			mockBlocks[0].TxHash,
			mockBlocks[0].ReceiptHash,
			mockBlocks[0].ExtraData,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mDB.sqlmock.ExpectCommit()

	if err = InsertBlock(mDB.dbConn, mockBlocks[0]); err != nil {
		t.Fatalf("failed to insert block: %v", err)
	}

	err = mDB.sqlmock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mDB.sqlmock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blocks" WHERE number = $1 ORDER BY "blocks"."hash" LIMIT $2`)).
		WithArgs(mockBlocks[0].Number, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"hash", "number", "gas_limit", "gas_used", "difficulty", "time",
			"parent_hash", "nonce", "miner", "size", "root_hash", "uncle_hash",
			"tx_hash", "receipt_hash", "extra_data",
		}).AddRow(
			mockBlocks[0].Hash, mockBlocks[0].Number, mockBlocks[0].GasLimit, mockBlocks[0].GasUsed, mockBlocks[0].Difficulty,
			mockBlocks[0].Time, mockBlocks[0].ParentHash, mockBlocks[0].Nonce, mockBlocks[0].Miner, mockBlocks[0].Size,
			mockBlocks[0].RootHash, mockBlocks[0].UncleHash, mockBlocks[0].TxHash, mockBlocks[0].ReceiptHash, mockBlocks[0].ExtraData,
		))

	retrievedBlock, err := GetBlockByNumber(mDB.dbConn, mockBlocks[0].Number)
	if err != nil {
		t.Fatalf("failed to get block by number: %v", err)
	}

	if !reflect.DeepEqual(retrievedBlock, &mockBlocks[0]) {
		t.Errorf("expected block:\n%+v\nbut got:\n%+v", mockBlocks[0], retrievedBlock)
	}

	err = mDB.sqlmock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("failed to meet expectations, got error: %v", err)
	}
}

func TestGetFirstBlock(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Set up the expectation for the SQL query
	mDB.sqlmock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blocks" ORDER BY number asc,"blocks"."hash" LIMIT $1`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"hash", "number", "gas_limit", "gas_used", "difficulty", "time",
			"parent_hash", "nonce", "miner", "size", "root_hash", "uncle_hash",
			"tx_hash", "receipt_hash", "extra_data",
		}).AddRow(
			mockBlocks[0].Hash, mockBlocks[0].Number, mockBlocks[0].GasLimit, mockBlocks[0].GasUsed, mockBlocks[0].Difficulty,
			mockBlocks[0].Time, mockBlocks[0].ParentHash, mockBlocks[0].Nonce, mockBlocks[0].Miner, mockBlocks[0].Size,
			mockBlocks[0].RootHash, mockBlocks[0].UncleHash, mockBlocks[0].TxHash, mockBlocks[0].ReceiptHash, mockBlocks[0].ExtraData,
		))

	retrievedBlock, err := GetFirstBlock(mDB.dbConn)
	if err != nil {
		t.Fatalf("failed to get first block: %v", err)
	}

	if !reflect.DeepEqual(retrievedBlock, &mockBlocks[0]) {
		t.Errorf("expected block:\n%+v\nbut got:\n%+v", mockBlocks[0], retrievedBlock)
	}

	err = mDB.sqlmock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("failed to meet expectations, got error: %v", err)
	}
}

func TestGetBlocks(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mDB.sqlmock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blocks"`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"hash", "number", "gas_limit", "gas_used", "difficulty", "time",
			"parent_hash", "nonce", "miner", "size", "root_hash", "uncle_hash",
			"tx_hash", "receipt_hash", "extra_data",
		}).AddRow(
			mockBlocks[0].Hash, mockBlocks[0].Number, mockBlocks[0].GasLimit, mockBlocks[0].GasUsed, mockBlocks[0].Difficulty,
			mockBlocks[0].Time, mockBlocks[0].ParentHash, mockBlocks[0].Nonce, mockBlocks[0].Miner, mockBlocks[0].Size,
			mockBlocks[0].RootHash, mockBlocks[0].UncleHash, mockBlocks[0].TxHash, mockBlocks[0].ReceiptHash, mockBlocks[0].ExtraData,
		).AddRow(
			mockBlocks[1].Hash, mockBlocks[1].Number, mockBlocks[1].GasLimit, mockBlocks[1].GasUsed, mockBlocks[1].Difficulty,
			mockBlocks[1].Time, mockBlocks[1].ParentHash, mockBlocks[1].Nonce, mockBlocks[1].Miner, mockBlocks[1].Size,
			mockBlocks[1].RootHash, mockBlocks[1].UncleHash, mockBlocks[1].TxHash, mockBlocks[1].ReceiptHash, mockBlocks[1].ExtraData,
		))

	retrievedBlocks, err := GetBlocks(mDB.dbConn)
	if err != nil {
		t.Fatalf("failed to get blocks: %v", err)
	}

	if len(retrievedBlocks) != len(mockBlocks) {
		t.Errorf("expected %d blocks but got %d", len(mockBlocks), len(retrievedBlocks))
	}

	for i, retrievedBlock := range retrievedBlocks {
		if !reflect.DeepEqual(retrievedBlock, &mockBlocks[i]) {
			t.Errorf("expected block:\n%+v\nbut got:\n%+v", mockBlocks[i], retrievedBlock)
		}
	}

	err = mDB.sqlmock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("failed to meet expectations, got error: %v", err)
	}
}
