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

	// 1. 读取address 的token balance
	output, err := client.ReadContract(usdtToken, "balanceOf(address)", addr.Hex())
	assert.NoError(t, err)
	rr, err := hexutil.DecodeBig(utls.FormatHex(hexutil.Encode(output)))
	assert.NoError(t, err)
	t.Log(rr.String()) // 地址的余额实时变化，这里就直接打印出来

	// 2. 读取合约的代币name()
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

	// 3. 读取合约的symbol
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

	// 4. 读取合约的total supply
	output, err = client.ReadContract(usdtToken, "totalSupply()")
	assert.NoError(t, err)

	in := utls.FormatHex(hexutil.Encode(output))
	total, err := hexutil.DecodeBig(in)
	assert.NoError(t, err)
	t.Log("total: ", total.String()) // usdt 总供应量会变，所以这里打印一下
	assert.Equal(t, "11377080729772723", total.String())

	// 5. 读取合约的owner
	output, err = client.ReadContract(usdtToken, "owner()")
	assert.NoError(t, err)
	t.Log(common.BytesToAddress(output).String()) // owner 也可能会被转移，还是打印一下
	assert.Equal(t, "0xC6CDE7C39eB2f0F0095F41570af89eFC2C1Ea828", common.BytesToAddress(output).String())

	// 6. 读取合约的 token decimals
	output, err = client.ReadContract(usdtToken, "decimals()")
	assert.NoError(t, err)
	hexDecimal := utls.FormatHex(hexutil.Encode(output))
	decimal, err := hexutil.DecodeUint64(hexDecimal)
	assert.NoError(t, err)
	assert.Equal(t, 6, decimal)
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
	tx, err := client.SendErc20TokenTransaction(false, senderPrv, nonce, gasLimit, gasPrice, reciver, tokenAddress, tokenAmount)
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
	tx, err := client.SendErc20TokenTransaction(true, senderPrv, nonce, gasLimit, gasPrice, reciver, tokenAddress, tokenAmount)
	assert.NoError(t, err)
	t.Log(tx.Hash().String())
}

// tool, 查询交易用
func TestERC20_SendErc20TokenTransaction(t *testing.T) {
	url := "https://rinkeby.infura.io/v3/36b98a13557c4b8583d57934ede2f74d"
	client, err := New(url)
	defer client.client.Close()
	assert.NoError(t, err)

	tx, ispending, err := client.client.TransactionByHash(context.Background(), common.HexToHash("0xe6e15d08b007f5ebc79fa1ae42201896cce9c39045d9ef42c54e3c8466e1e990"))
	assert.NoError(t, err)
	t.Log(ispending)
	t.Log(hexutil.Encode(tx.Data()))

}
