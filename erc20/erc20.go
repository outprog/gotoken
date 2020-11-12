package erc20

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zyjblockchain/sandy_log/log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20 struct {
	ChainId *big.Int
	client  *ethclient.Client
}

func New(web3 string) (*ERC20, error) {
	client, err := ethclient.Dial(web3)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	return &ERC20{
		ChainId: chainId,
		client:  client,
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

func (e *ERC20) Approve() {

}

// ReadContract 读取合约通用接口
// arg 为hex 类型，如果有多个arg，顺序必须和funcName 中的顺序一致
// 返回链上原始的 output bytes
func (e *ERC20) ReadContract(contractAddress common.Address, funcName string, arg ...string) ([]byte, error) {
	funcCode := getContractFunctionCode(funcName)
	argCode := make([]byte, 0)
	if len(arg) > 0 {
		for _, val := range arg {
			argCode = append(argCode, formatArgs(val)...)
		}
	}
	inputData := make([]byte, 0)
	inputData = append(funcCode, argCode...)

	// call
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: inputData,
	}
	output, err := e.client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	log.Infof("raw output: %s", hexutil.Encode(output))
	return output, nil
}

// SendErc20TokenTransaction 转账交易或者是approve 交易
func (e *ERC20) SendErc20TokenTransaction(isApprove bool, private string, nonce, gasLimit uint64, gasPrice *big.Int, receiver, tokenAddress common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	rawTx := newErc20TokenRawTx(isApprove, nonce, receiver, tokenAddress, gasLimit, gasPrice, tokenAmount)

	// 对原生交易进行签名
	prv, err := crypto.ToECDSA(common.FromHex(private))
	if err != nil {
		log.Errorf("crypto.ToECDSA error: %v", err)
		return nil, err
	}
	signedTx, err := signRawTx(rawTx, e.ChainId, prv)
	if err != nil {
		log.Errorf("signer transaction error: %v, rawTx: %v", err, rawTx)
		return nil, err
	}

	// 发送签好名的交易上链
	err = e.client.SendTransaction(context.Background(), signedTx)
	// todo 业务需求可以自己处理error
	return signedTx, err
}
