package erc20

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20 struct {
	client *ethclient.Client
}

func New(web3 string) (*ERC20, error) {
	client, err := ethclient.Dial(web3)
	if err != nil {
		return nil, err
	}

	return &ERC20{
		client: client,
	}, nil
}

func (e *ERC20) BalanceOf(token, owner common.Address) (*big.Int, error) {
	data, err := BalanceOf(owner)
	if err != nil {
		return nil, err
	}

	res, err := e.client.CallContract(context.TODO(), ethereum.CallMsg{
		Data: data,
		To:   &token,
	}, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(res), nil
}
