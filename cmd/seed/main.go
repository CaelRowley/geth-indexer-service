package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/CaelRowley/geth-indexer-service/pkg/data"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	BlockCount = 500
	TxCount    = 2000
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "err", err)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")),
		&gorm.Config{
			SkipDefaultTransaction: true,
		})
	if err != nil {
		slog.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}

	if err := db.Exec("DELETE FROM blocks; DELETE FROM transactions;").Error; err != nil {
		slog.Error("failed to clear tables", "err", err)
		os.Exit(1)
	}
	slog.Info("db cleared")
	if err := db.Create(createBlocks()).Error; err != nil {
		slog.Error("failed to seed blocks", "err", err)
		os.Exit(1)
	}
	slog.Info("blocks seeded")
	if err := db.Create(createTxs()).Error; err != nil {
		slog.Error("failed to seed txs", "err", err)
		os.Exit(1)
	}
	slog.Info("transactions seeded")
}

func createBlocks() []data.Block {
	var blocks []data.Block
	for i := 0; i < 500; i++ {
		blocks = append(blocks, data.Block{
			Hash:        "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			Number:      uint64(i),
			GasLimit:    8000000 + uint64(i),
			GasUsed:     7500000 + uint64(i),
			Difficulty:  1000000 + uint64(i),
			Time:        1627891200 + uint64(i),
			ParentHash:  "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			Nonce:       0,
			Miner:       "0x000000000000000000000000000000000000" + fmt.Sprintf("%04d", i),
			Size:        1000 + uint64(i),
			RootHash:    "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UncleHash:   "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			TxHash:      "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			ReceiptHash: "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			ExtraData:   []byte("extra data " + strconv.Itoa(i)),
		})
	}
	return blocks
}

func createTxs() []data.Transaction {
	var txs []data.Transaction
	for i := 0; i < TxCount; i++ {
		newTx := data.Transaction{
			Hash:      "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			From:      "0x000000000000000000000000000000000000" + fmt.Sprintf("%04d", i),
			Contract:  "0x" + fmt.Sprintf("%04d", i) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			Value:     200 + uint64(i),
			Data:      []byte("some data " + strconv.Itoa(i)),
			Gas:       500000 + uint64(i),
			GasPrice:  1000000 + uint64(i),
			Cost:      1000000000 + uint64(i),
			Nonce:     0,
			Status:    uint64(i),
			BlockHash: "0x" + fmt.Sprintf("%04d", i%BlockCount) + "efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		}
		if i%2 == 0 {
			newTx.To = "0x000000000000000000000000000000000000" + fmt.Sprintf("%04d", i+1)
		}
		txs = append(txs, newTx)
	}
	return txs
}
