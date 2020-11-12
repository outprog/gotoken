package erc20

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/outprog/gotoken/utls"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
	"unicode"
)

func TestERC20(t *testing.T) {
	url := "https://mainnet.infura.io/v3/f1301efa5af1432c84063f231f08f920"
	usdtToken := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7") // usdt
	addr := common.HexToAddress("0x1C76Caf4D8ADeD815f76C958a70702fe92aC4982")
	client, err := New(url)
	assert.NoError(t, err)
	bal, err := client.BalanceOf(usdtToken, addr)
	assert.NoError(t, err)
	t.Log(bal.String())
}

func TestERC20_ReadContract(t *testing.T) {
	url := "https://mainnet.infura.io/v3/f1301efa5af1432c84063f231f08f920"
	usdtToken := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7") // usdt
	addr := common.HexToAddress("0x1C76Caf4D8ADeD815f76C958a70702fe92aC4982")
	client, err := New(url)
	assert.NoError(t, err)

	// 1. 读取address 的token balance: `function balanceOf(address _owner) constant returns (uint256 balance)`
	output, err := client.ReadContract(usdtToken, "balanceOf(address)", addr.Hex())
	assert.NoError(t, err)
	rr, err := hexutil.DecodeBig(utls.FormatHex(hexutil.Encode(output)))
	assert.NoError(t, err)
	t.Log(rr.String()) // 地址的余额实时变化，这里就直接打印出来

	// 2. 读取合约的代币name(): `function name() constant returns (string name) `
	output, err = client.ReadContract(usdtToken, "name()")
	assert.NoError(t, err)
	// 处理返回的数据
	aa := string(output)
	dd := make([]rune, 0)
	for _, val := range aa {
		if unicode.IsGraphic(val) {
			dd = append(dd, val)
		}
	}
	assert.Equal(t, "Tether USD", strings.TrimSpace(string(dd)))

	// 3. 读取合约的symbol: `function symbol() constant returns (string symbol)`
	output, err = client.ReadContract(usdtToken, "symbol()")
	assert.NoError(t, err)
	// 处理返回的数据
	aa = string(output)
	dd = make([]rune, 0)
	for _, val := range aa {
		if unicode.IsGraphic(val) {
			dd = append(dd, val)
		}
	}
	assert.Equal(t, "USDT", strings.TrimSpace(string(dd)))

	// 4. 读取合约的total supply: `function totalSupply() constant returns (uint256 totalSupply)`
	output, err = client.ReadContract(usdtToken, "totalSupply()")
	assert.NoError(t, err)

	in := utls.FormatHex(hexutil.Encode(output))
	total, err := hexutil.DecodeBig(in)
	assert.NoError(t, err)
	t.Log("total: ", total.String()) // usdt 总供应量会变，所以这里打印一下
	assert.Equal(t, "11377080729772723", total.String())

	// 5. 读取合约的owner: `function owner() constant returns (string owner)` 这个不是erc20 标准中的固定方法
	output, err = client.ReadContract(usdtToken, "owner()")
	assert.NoError(t, err)
	t.Log(common.BytesToAddress(output).String()) // owner 也可能会被转移，还是打印一下
	assert.Equal(t, "0xC6CDE7C39eB2f0F0095F41570af89eFC2C1Ea828", common.BytesToAddress(output).String())

	// 6. 读取合约的 token decimals: `function decimals() constant returns (uint8 decimals)`
	output, err = client.ReadContract(usdtToken, "decimals()")
	assert.NoError(t, err)
	hexDecimal := utls.FormatHex(hexutil.Encode(output))
	decimal, err := hexutil.DecodeUint64(hexDecimal)
	assert.NoError(t, err)
	assert.Equal(t, uint64(6), decimal)

	// 7. 读取合约的 token allowance: `function allowance(address _owner, address _spender) constant returns (uint256 remaining)`
	url = "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err = New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	tokenowner := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	tokenspender := "0xF9891E1A2635CB8D8C25A6A2ec8E453bFb2E67c4"
	tokenAddr := common.HexToAddress("0x03332638A6b4F5442E85d6e6aDF929Cd678914f1")
	output, err = client.ReadContract(tokenAddr, "allowance(address,address)", tokenowner, tokenspender)
	assert.NoError(t, err)
	rr, err = hexutil.DecodeBig(utls.FormatHex(hexutil.Encode(output)))
	assert.NoError(t, err)
	t.Log(rr.String())
}

func TestERC20_SendErc20TokenTransaction_Transfer(t *testing.T) {
	url := "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err := New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	sender := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	senderPrv := "69f657eaf364969ccfb2531f45d9c9efac0a63e359cea51e5f7d8340784168d2"   // 用户测试私钥，请勿乱用
	tokenAddress := common.HexToAddress("0x03332638A6b4F5442E85d6e6aDF929Cd678914f1") // Test3 token address
	nonce, err := client.client.NonceAt(context.Background(), common.HexToAddress(sender), nil)
	assert.NoError(t, err)

	tokenAmount, _ := new(big.Int).SetString("888888888888888888", 10) // 测试的这个token 的decimals 为 18位，所以需要大一点
	reciver := common.HexToAddress("0xF9891E1A2635CB8D8C25A6A2ec8E453bFb2E67c4")

	gasLimit := uint64(60000)
	gasPrice, _ := client.client.SuggestGasPrice(context.Background())
	// token 转账交易
	// 1. 构造raw transaction
	rawTx := NewErc20TokenTransferOrApproveRawTx(false, nonce, reciver, tokenAddress, gasLimit, gasPrice, tokenAmount)
	// 2. 签名并发送交易
	tx, err := client.SignAndSendTransaction(senderPrv, rawTx)
	assert.NoError(t, err)
	t.Log(tx.Hash().String())
}

func TestERC20_SendErc20TokenTransaction_Approve(t *testing.T) {
	url := "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err := New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	sender := "0x59375A522876aB96B0ed2953D0D3b92674701Cc2"
	senderPrv := "69f657eaf364969ccfb2531f45d9c9efac0a63e359cea51e5f7d8340784168d2"   // 用户测试私钥，请勿乱用
	tokenAddress := common.HexToAddress("0x03332638A6b4F5442E85d6e6aDF929Cd678914f1") // Test3 token address
	nonce, err := client.client.NonceAt(context.Background(), common.HexToAddress(sender), nil)
	assert.NoError(t, err)

	tokenAmount, _ := new(big.Int).SetString("9999999999999999999", 10) // 测试的这个token 的decimals 为 18位，所以需要大一点
	reciver := common.HexToAddress("0xF9891E1A2635CB8D8C25A6A2ec8E453bFb2E67c4")

	gasLimit := uint64(60000)
	gasPrice, _ := client.client.SuggestGasPrice(context.Background())
	// token approve 交易
	// 1. 构造raw transaction
	rawTx := NewErc20TokenTransferOrApproveRawTx(true, nonce, reciver, tokenAddress, gasLimit, gasPrice, tokenAmount)
	// 2. 签名并发送
	tx, err := client.SignAndSendTransaction(senderPrv, rawTx)
	assert.NoError(t, err)
	t.Log(tx.Hash().String())
}

func TestNewErc20TokenTransferFromRawTx(t *testing.T) {
	url := "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err := New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	sender := "0xF9891E1A2635CB8D8C25A6A2ec8E453bFb2E67c4"
	senderPrv := "701081097d795a34f59b5b2795938057f879704e81f26de833f5199b3256f709" // 不要乱用
	nonce, err := client.client.NonceAt(context.Background(), common.HexToAddress(sender), nil)
	assert.NoError(t, err)

	tokenOwner := common.HexToAddress("0x59375A522876aB96B0ed2953D0D3b92674701Cc2")
	tokenReceiver := common.HexToAddress("0x811Ab218e53E2125d9311650a448782737Fb6E42")
	tokenAddress := common.HexToAddress("0x03332638A6b4F5442E85d6e6aDF929Cd678914f1") // Test3 token address
	tokenAmount, _ := new(big.Int).SetString("1111111111111111111", 10)

	gasLimit := uint64(100000) // 注：这里gasLimit 会大于60000
	gasPrice, _ := client.client.SuggestGasPrice(context.Background())

	// transferFrom 交易
	// 1. 构造raw transaction
	rawTx := NewErc20TokenTransferFromRawTx(nonce, tokenOwner, tokenReceiver, tokenAddress, tokenAmount, gasLimit, gasPrice)
	// 2. 签名并发送
	tx, err := client.SignAndSendTransaction(senderPrv, rawTx)
	assert.NoError(t, err)
	t.Log(tx.Hash().String())
}

// tool, 查询交易用
func TestERC20_SendErc20TokenTransaction(t *testing.T) {
	url := "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err := New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	txHash := common.HexToHash("0x7f9aacb1c7359770bfc81aaa1f5858fc0042a79f9c61aff6d498b2a40a7ea0fe")
	tx, ispending, err := client.client.TransactionByHash(context.Background(), txHash)
	assert.NoError(t, err)
	t.Log(ispending)
	t.Log(hexutil.Encode(tx.Data()))

}
