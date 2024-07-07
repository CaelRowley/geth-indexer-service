package data

type Transaction struct {
	Hash      string `json:"hash" gorm:"column:hash;type:char(66);primaryKey"`
	From      string `json:"from" gorm:"column:from;type:char(42);not null"`
	To        string `json:"to" gorm:"column:to;type:char(42)"`
	Contract  string `json:"contract" gorm:"column:contract;type:char(66);not null"`
	Value     uint64 `json:"value" gorm:"column:value;type:numeric;not null"`
	Data      []byte `json:"data" gorm:"column:data;type:bytea;not null"`
	Gas       uint64 `json:"gas" gorm:"column:gas;type:numeric;not null"`
	GasPrice  uint64 `json:"gasPrice" gorm:"column:gas_price;type:numeric;not null"`
	Cost      uint64 `json:"cost" gorm:"column:cost;type:numeric;not null"`
	Nonce     uint64 `json:"nonce" gorm:"column:nonce;type:numeric;not null"`
	Status    uint64 `json:"status" gorm:"column:status;type:numeric;not null"`
	BlockHash string `json:"blockHash" gorm:"column:block_hash;type:char(66);not null"`
}
