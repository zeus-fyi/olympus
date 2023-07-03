package periphery

import (
	_ "embed"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
)

//go:embed contracts/interfaces/IPeripheryPaymentsWithFee.sol/IPeripheryPaymentsWithFee.json
var paymentsABI []byte

type FeeOptions struct {
	Fee       *entities.Percent // The percent of the output that will be taken as a fee.
	Recipient accounts.Address  // The recipient of the fee.

}

func encodeFeeBips(fee *entities.Percent) *big.Int {
	return fee.Multiply(entities.NewPercent(big.NewInt(10000), big.NewInt(1))).Quotient()
}

func EncodeUnwrapWETH9(amountMinimum *big.Int, recipient accounts.Address, feeOptions *FeeOptions) ([]byte, error) {
	abi := GetABI(paymentsABI)
	if feeOptions != nil {
		return abi.Pack("unwrapWETH9WithFee", amountMinimum, &recipient, encodeFeeBips(feeOptions.Fee), feeOptions.Recipient)
	}

	return abi.Pack("unwrapWETH9", amountMinimum, recipient)
}

func EncodeSweepToken(token *entities.Token, amountMinimum *big.Int, recipient accounts.Address, feeOptions *FeeOptions) ([]byte, error) {
	abi := GetABI(paymentsABI)

	if feeOptions != nil {
		return abi.Pack("sweepTokenWithFee", token.Address, amountMinimum, recipient, encodeFeeBips(feeOptions.Fee), feeOptions.Recipient)
	}

	return abi.Pack("sweepToken", token.Address, amountMinimum, recipient)
}

func EncodeRefundETH() []byte {
	abi := GetABI(paymentsABI)
	data, err := abi.Pack("refundETH")
	if err != nil {
		panic(err)
	}
	return data
}
