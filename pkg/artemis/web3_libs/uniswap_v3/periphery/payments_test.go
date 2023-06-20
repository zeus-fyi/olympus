package periphery

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeus-fyi/gochain/web3/accounts"
	entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
)

var (
	recipient  = accounts.HexToAddress("0x0000000000000000000000000000000000000003")
	amount     = big.NewInt(123)
	feeOptions = &FeeOptions{
		Fee:       entities.NewPercent(big.NewInt(1), big.NewInt(1000)),
		Recipient: accounts.HexToAddress("0x0000000000000000000000000000000000000009"),
	}
	token = entities.NewToken(1, accounts.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "token0")
)

func TestEncodeUnwrapWETH9(t *testing.T) {
	// works without feeOptions
	calldata, err := EncodeUnwrapWETH9(amount, recipient, nil)
	assert.NoError(t, err)
	assert.Equal(t, "0x49404b7c000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000003", hexutil.Encode(calldata))

	// works with feeOptions
	calldata, err = EncodeUnwrapWETH9(amount, recipient, feeOptions)
	assert.NoError(t, err)
	assert.Equal(t, "0x9b2c0a37000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000009", hexutil.Encode(calldata))
}

func TestEncodeSweepToken(t *testing.T) {
	// works without feeOptions
	calldata, err := EncodeSweepToken(token, amount, recipient, nil)
	assert.NoError(t, err)
	assert.Equal(t, "0xdf2ab5bb0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000003", hexutil.Encode(calldata))

	// works with feeOptions
	calldata, err = EncodeSweepToken(token, amount, recipient, feeOptions)
	assert.NoError(t, err)
	assert.Equal(t, "0xe0e189a00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000009", hexutil.Encode(calldata))
}

func TestEncodeRefundETH(t *testing.T) {
	// works without feeOptions
	calldata := EncodeRefundETH()
	assert.Equal(t, "0x12210e8a", hexutil.Encode(calldata))
}
