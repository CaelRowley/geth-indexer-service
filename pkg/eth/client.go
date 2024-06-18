package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	return client, nil
}
