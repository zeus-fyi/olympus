package async_analysis

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
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

func (c *ContractAnalysis) CalculateTransferFeeTax(ctx context.Context, amount *big.Int) (*uniswap_core_entities.Percent, error) {
	err := c.UserA.SetERC20BalanceBruteForce(ctx, c.SmartContractAddr, c.UserA.Address().String(), amount)
	if err != nil {
		return nil, err
	}
	startBalUserB, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserB.Address().String())
	if err != nil {
		return nil, err
	}
	initBalUserA, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserA.Address().String())
	if err != nil {
		return nil, err
	}
	if initBalUserA.String() != amount.String() {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("erc20Diff is not equal to amount")
		return uniswap_core_entities.NewPercent(big.NewInt(0), big.NewInt(1)), errors.New("erc20Diff is 0")
	}
	c.MethodName = "transfer"
	c.Params = []interface{}{c.UserB.Address(), amount}
	tx, err := c.UserA.TransferERC20Token(ctx, c.SendContractTxPayload)
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	fmt.Println("tx hash", tx.Hash().String())
	rx, err := c.UserA.GetTxReceipt(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}
	if rx.Status == types.ReceiptStatusFailed {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("tx failed, amount is 0")
		return uniswap_core_entities.NewPercent(big.NewInt(0), big.NewInt(1)), errors.New("tx failed, amount is 0")
	}
	endBalUserB, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserB.Address().String())
	if err != nil {
		return nil, err
	}
	if endBalUserB.String() == "0" {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("transfer amount is 0")
		return uniswap_core_entities.NewPercent(big.NewInt(0), big.NewInt(1)), errors.New("transfer amount is 0")
	}
	transferAmount := new(big.Int).Sub(endBalUserB, startBalUserB)
	feeAmount := new(big.Int).Sub(amount, transferAmount)
	total := new(big.Int).Add(transferAmount, feeAmount)
	if total.String() != amount.String() {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("total is not equal to amount")
		return uniswap_core_entities.NewPercent(big.NewInt(0), big.NewInt(1)), errors.New("total is not equal to amount")
	}
	if feeAmount.String() == "0" {
		return uniswap_core_entities.NewPercent(big.NewInt(1), big.NewInt(1)), nil
	}
	gcd := new(big.Int).GCD(nil, nil, amount, feeAmount)
	numerator := new(big.Int).Div(feeAmount, gcd)
	denominator := new(big.Int).Div(amount, gcd)
	percent := uniswap_core_entities.NewPercent(numerator, denominator)
	return percent, nil
}
