# Read and perform contract (power of golang)

### 已完成
1. 读取合约中的所有constant功能
2. 执行合约方法：transfer(), approve(),transferFrom()

---
### Example 
```go
    // 1. 创建client
    nodeUrl := "https://mainnet.infura.io/v3/f1301efa5af1432c84063f231f08f920" // 测试可以，别乱用
    client, err := New(nodeUrl)

    // 2. 读取合约
    output, err := client.ReadContract(usdtToken, "balanceOf(address)", addr.Hex())
    assert.NoError(t, err)
    rr, err := hexutil.DecodeBig(utls.FormatHex(hexutil.Encode(output)))
    assert.NoError(t, err)
    t.Log(rr.String())

    // 3. 执行合约方法
    // token 转账交易
	// 1. 构造raw transaction
	rawTx := NewErc20TokenTransferOrApproveRawTx(false, nonce, reciver, tokenAddress, gasLimit, gasPrice, tokenAmount)
	// 2. 签名并发送交易
	tx, err := client.SignAndSendTransaction(senderPrv, rawTx)

```
#### 注：在 `erc20/erc20_test.go` 中可以看到erc20 合约的每一个方法的调用测试