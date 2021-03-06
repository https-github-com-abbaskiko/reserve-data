package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-data/common/blockchain"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// GeneratedSetBaseRate build tx set base rate
func (bc *Blockchain) GeneratedSetBaseRate(opts blockchain.TxOpts, tokens []ethereum.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return bc.BuildTx(timeout, opts, bc.pricing, "setBaseRate", tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

// GeneratedSetCompactData build tx to set compact data
func (bc *Blockchain) GeneratedSetCompactData(opts blockchain.TxOpts, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return bc.BuildTx(timeout, opts, bc.pricing, "setCompactData", buy, sell, blockNumber, indices)
}

// GeneratedGetRate get token rate from reserve
func (bc *Blockchain) GeneratedGetRate(opts blockchain.CallOpts, token ethereum.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	timeOut := 2 * time.Second
	out := big.NewInt(0)
	err := bc.Call(timeOut, opts, bc.pricing, out, "getRate", token, currentBlockNumber, buy, qty)
	return out, err
}

// GeneratedGetListedTokens return listed tokens on reserve
func (bc *Blockchain) GeneratedGetListedTokens(opts blockchain.CallOpts) ([]ethereum.Address, error) {
	var (
		ret0 = new([]ethereum.Address)
	)
	timeout := 2 * time.Second
	out := ret0
	err := bc.Call(timeout, opts, bc.pricing, out, "getListedTokens")
	return *ret0, err
}
