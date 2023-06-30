package async_analysis

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type ContractAnalysis struct {
	u *web3_client.UniswapClient
	web3_actions.SendContractTxPayload
	UserA web3_client.Web3Client
	UserB web3_client.Web3Client
}

func NewERC20ContractAnalysis(u *web3_client.UniswapClient, address string) ContractAnalysis {
	return NewContractAnalysis(u, address, web3_client.MustLoadERC20Abi())
}

func NewContractAnalysis(u *web3_client.UniswapClient, address string, abiFile *abi.ABI) ContractAnalysis {
	ca := ContractAnalysis{
		u: u,
		SendContractTxPayload: web3_actions.SendContractTxPayload{
			SmartContractAddr: address,
			ContractABI:       abiFile,
		},
	}
	ca.UserA = u.Web3Client
	return ca
}

func (c *ContractAnalysis) CalculateTransferFeeTax(ctx context.Context, amount *big.Int) (int64, error) {
	err := c.u.Web3Client.SetERC20BalanceBruteForce(ctx, c.SmartContractAddr, c.UserA.Address().String(), amount)
	if err != nil {
		return -1, err
	}
	startBal, err := c.u.Web3Client.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserB.Address().String())
	if err != nil {
		return -1, err
	}
	c.MethodName = "transfer"
	c.Params = []interface{}{c.UserB.Address(), amount}
	_, err = c.u.Web3Client.TransferERC20Token(ctx, c.SendContractTxPayload)
	if err != nil {
		return -1, err
	}
	endBal, err := c.u.Web3Client.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserB.Address().String())
	if err != nil {
		return -1, err
	}
	transferAmount := new(big.Int).Sub(endBal, startBal)
	feeAmount := new(big.Int).Sub(amount, transferAmount)
	fmt.Println("amount", amount.String())
	fmt.Println("transferAmount", transferAmount.String())
	fmt.Println("feeAmount", feeAmount.String())
	transferFeePercent := new(big.Int).Div(new(big.Int).Mul(feeAmount, new(big.Int).SetUint64(100)), amount)
	return transferFeePercent.Int64(), nil
}
