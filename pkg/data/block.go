package data

type Block struct {
	Hash        string `json:"hash"`
	Number      uint64 `json:"number"`
	GasLimit    uint64 `json:"gasLimit"`
	GasUsed     uint64 `json:"gasUsed"`
	Difficulty  string `json:"difficulty"`
	Time        uint64 `json:"time"`
	ParentHash  string `json:"parentHash"`
	Nonce       string `json:"nonce"`
	Miner       string `json:"miner"`
	Size        uint64 `json:"size"`
	RootHash    string `json:"rootHash"`
	UncleHash   string `json:"uncleHash"`
	TxHash      string `json:"txHash"`
	ReceiptHash string `json:"receiptHash"`
	ExtraData   []byte `json:"extraData"`
}
