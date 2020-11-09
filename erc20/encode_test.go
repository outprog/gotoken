package erc20

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestBalanceOf(t *testing.T) {
	balABI, err := BalanceOf(common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"))
	assert.NoError(t, err)
	assert.Equal(t, "0x70a08231000000000000000000000000a06b79e655db7d7c3b3e7b2cceeb068c3259d0c9", common.ToHex(balABI))
}
