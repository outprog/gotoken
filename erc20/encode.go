package erc20

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var erc20ABI abi.ABI

func init() {
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(Erc20ABI))
	if err != nil {
		panic(err)
	}
}

func Approve(spender common.Address, value *big.Int) ([]byte, error) {
	return erc20ABI.Pack("approve", spender, value)
}

func BalanceOf(owner common.Address) ([]byte, error) {
	return erc20ABI.Pack("balanceOf", owner)
}

// getContractFunctionCode 计算合约函数code
func getContractFunctionCode(funcName string) []byte {
	h := crypto.Keccak256Hash([]byte(funcName))
	return h.Bytes()[:4]
}

// formatArgs 把参数转换成[32]byte的数组类型
func formatArgs(args string) []byte {
	b := common.FromHex(args)
	var h [32]byte
	if len(b) > len(h) {
		b = b[len(b)-32:]
	}
	copy(h[32-len(b):], b)
	return h[:]
}

// NewErc20TokenTransferOrApproveRawTx 构造erc20 token tansfer 或者 approve 的 raw transaction
func NewErc20TokenTransferOrApproveRawTx(approveTx bool, senderNonce uint64, receiver common.Address, contractAddr common.Address, gasLimit uint64, gasPrice *big.Int, tokenAmount *big.Int) *types.Transaction {
	/**
	transferFun := "0xa9059cbb"
	receiverAddrCode := 000000000000000000000000b1e15fdbe88b7e7c47552e2d33cd5a9b2e0fd478 // eg: 代币接收地址code
	tokenAmountCode := "0000000000000000000000000000000000000000000000000000000000000064" // eg: 转币数量100
	*/
	funcName := "transfer(address,uint256)"
	if approveTx {
		funcName = "approve(address,uint256)"
	}
	funcCode := getContractFunctionCode(funcName)
	receiverAddrCode := formatArgs(receiver.Hex())
	AmountCode := formatArgs(tokenAmount.Text(16)) // big.Int 转 hex

	// 组合生成执行合约的input
	inputData := make([]byte, 0)
	inputData = append(append(funcCode, receiverAddrCode...), AmountCode...) // 顺序千万不能乱，可以在etherscan上找个合约交易查看input data

	// 组装以太坊交易
	return types.NewTransaction(senderNonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, inputData)
}

// NewErc20TokenTransferFromRawTx
// 在approve 之后使用， tokenOwner为授权token 出去的地址，这笔交易的签名者为授权token 的使用者，tokenReceiver 为这笔转账的token 接收者；
func NewErc20TokenTransferFromRawTx(senderNonce uint64, tokenOwner, tokenReceiver, tokenAddr common.Address, tokenAmount *big.Int, gasLimit uint64, gasPrice *big.Int) *types.Transaction {
	funcName := "transferFrom(address,address,uint256)" // transferFrom(address _from, address _to, uint256 _value)

	funcCode := getContractFunctionCode(funcName)
	tokenOwnerCode := formatArgs(tokenOwner.Hex())
	tokenReceiverCode := formatArgs(tokenReceiver.Hex())
	tokenAmountCode := formatArgs(tokenAmount.Text(16))

	// input
	inputData := make([]byte, 0)
	inputData = append(append(append(funcCode, tokenOwnerCode...), tokenReceiverCode...), tokenAmountCode...) // 顺序一致

	return types.NewTransaction(senderNonce, tokenAddr, big.NewInt(0), gasLimit, gasPrice, inputData)
}

// signRawTx 对交易进行签名
func signRawTx(rawTx *types.Transaction, chainID *big.Int, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(rawTx, signer, prv)
	return signedTx, err
}
