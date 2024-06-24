package data

type Block struct {
	Hash        string `json:"hash" gorm:"column:hash;type:char(66);primaryKey"`
	Number      uint64 `json:"number" gorm:"column:number;type:numeric;not null;unique;index:,sort:asc"`
	GasLimit    uint64 `json:"gasLimit" gorm:"column:gas_limit;type:numeric;not null"`
	GasUsed     uint64 `json:"gasUsed" gorm:"column:gas_used;type:numeric;not null"`
	Difficulty  string `json:"difficulty" gorm:"column:difficulty;type:varchar;not null"`
	Time        uint64 `json:"time" gorm:"column:time;type:numeric;not null"`
	ParentHash  string `json:"parentHash" gorm:"column:parent_hash;type:char(66);not null"`
	Nonce       string `json:"nonce" gorm:"column:nonce;type:varchar;not null"`
	Miner       string `json:"miner" gorm:"column:miner;type:char(42);not null"`
	Size        uint64 `json:"size" gorm:"column:size;type:numeric;not null"`
	RootHash    string `json:"rootHash" gorm:"column:root_hash;type:char(66);not null"`
	UncleHash   string `json:"uncleHash" gorm:"column:uncle_hash;type:char(66);not null"`
	TxHash      string `json:"txHash" gorm:"column:tx_hash;type:char(66);not null"`
	ReceiptHash string `json:"receiptHash" gorm:"column:receipt_hash;type:char(66);not null"`
	ExtraData   []byte `json:"extraData" gorm:"column:extra_data;type:bytea"`
}
