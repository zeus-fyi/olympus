package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (w *Web3Client) GetBlockTxs(ctx context.Context) (types.Transactions, error) {
	w.Dial()
	defer w.Close()
	block, err := w.C.BlockByNumber(ctx, nil)
	if err != nil {
		log.Err(err).Msg("failed to get nonce")
		return nil, err
	}
	return block.Transactions(), nil
}

func (w *Web3Client) GetTxByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	w.Dial()
	defer w.Close()
	tx, isPending, err := w.C.TransactionByHash(ctx, hash)
	if err != nil {
		log.Err(err).Msg("failed to get nonce")
		return nil, false, err
	}
	return tx, isPending, nil
}

/*
0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
239999999999999788
0xf9974c3357Ba17a6B8A8326B65233e17c83AFc91
*/

// Eth in -> WETH out -> token out

func (u *UniswapClient) ExecTradeV2SwapPayable(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
	// todo max this window more appropriate vs near infinite

	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{},
		Deadline: sigDeadline,
		Payable:  nil,
	}

	sc1 := UniversalRouterExecSubCmd{
		Command:   WrapETH,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: WrapETHParams{
			Recipient: accounts.HexToAddress(universalRouterRecipient),
			AmountMin: to.AmountIn,
		},
	}
	ur.Commands = append(ur.Commands, sc1)
	sc2 := UniversalRouterExecSubCmd{
		Command:   V2SwapExactIn,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: V2SwapExactInParams{
			AmountIn:      to.AmountIn,
			AmountOutMin:  to.AmountOut,
			Path:          []accounts.Address{to.AmountInAddr, to.AmountOutAddr},
			To:            accounts.HexToAddress(universalRouterSender),
			PayerIsSender: false,
		},
	}
	ur.Commands = append(ur.Commands, sc2)
	sc3 := UniversalRouterExecSubCmd{
		Command:   UnwrapWETH,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: UnwrapWETHParams{
			Recipient: accounts.HexToAddress(universalRouterSender),
			AmountMin: new(big.Int).SetUint64(0),
		},
	}
	ur.Commands = append(ur.Commands, sc3)
	payable := &web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    to.AmountIn,
			ToAddress: u.Web3Client.Address(),
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	ur.Payable = payable

	tx, err := u.ExecUniswapUniversalRouterCmd(ur)
	if err != nil {
		return err
	}
	to.AddTxHash(accounts.Hash(tx.Hash()))
	return err
}

func (u *UniswapClient) ExecTradeV2SwapFromToken(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
	// todo max this window more appropriate vs near infinite

	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{},
		Deadline: sigDeadline,
		Payable:  nil,
	}

	max, _ := new(big.Int).SetString(maxUINT, 10)
	approveTx, err := u.ApproveSpender(ctx, to.AmountInAddr.String(), Permit2SmartContractAddress, max)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	sc1 := UniversalRouterExecSubCmd{
		Command:   Permit2Permit,
		CanRevert: false,
		Inputs:    nil,
	}

	fmt.Println(to.AmountInAddr.String(), "to.AmountInAddr.String()")
	fmt.Println(to.AmountOutAddr.String(), "to.AmountOutAddr.String()")
	psp := Permit2PermitParams{
		PermitSingle{
			PermitDetails: PermitDetails{
				Token:      to.AmountInAddr,
				Amount:     to.AmountIn,
				Expiration: sigDeadline,
				// todo this needs to update a nonce count in db or track them somehow
				Nonce: new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress(UniswapUniversalRouterAddress),
			SigDeadline: sigDeadline,
		},
		nil,
	}
	err = psp.SignPermit2Mainnet(u.Web3Client.Account)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return err
	}
	if psp.Signature == nil {
		log.Warn().Msg("signature is nil")
		return errors.New("signature is nil")
	}
	sc1.DecodedInputs = psp
	ur.Commands = append(ur.Commands, sc1)
	sc2 := UniversalRouterExecSubCmd{
		Command:   V2SwapExactIn,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: V2SwapExactInParams{
			AmountIn:      to.AmountIn,
			AmountOutMin:  to.AmountOut,
			Path:          []accounts.Address{to.AmountInAddr, to.AmountOutAddr},
			To:            accounts.HexToAddress(u.Web3Client.Address().Hex()),
			PayerIsSender: true,
		},
	}
	ur.Commands = append(ur.Commands, sc2)
	sc3 := UniversalRouterExecSubCmd{
		Command:   UnwrapWETH,
		CanRevert: true,
		Inputs:    nil,
		DecodedInputs: UnwrapWETHParams{
			Recipient: accounts.HexToAddress(u.Web3Client.PublicKey()),
			AmountMin: new(big.Int).SetUint64(0),
		},
	}
	ur.Commands = append(ur.Commands, sc3)
	//payable := &web3_actions.SendEtherPayload{
	//	TransferArgs: web3_actions.TransferArgs{
	//		Amount:    to.AmountIn,
	//		ToAddress: u.Web3Client.Address(),
	//	},
	//	GasPriceLimits: web3_actions.GasPriceLimits{},
	//}
	//ur.Payable = payable

	tx, err := u.ExecUniswapUniversalRouterCmd(ur)
	if err != nil {
		return err
	}
	to.AddTxHash(accounts.Hash(tx.Hash()))
	return err
}

func (u *UniswapClient) RouterApproveAndSend(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
	approveTx, err := u.ApproveSpender(ctx, to.AmountInAddr.String(), u.RouterSmartContractAddr, to.AmountIn)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(accounts.Hash(approveTx.Hash()))
	transferTxParams := web3_actions.SendContractTxPayload{
		SmartContractAddr: to.AmountInAddr.String(),
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				ToAddress: accounts.HexToAddress(pairContractAddr),
			},
		},
		ContractABI: MustLoadERC20Abi(),
		Params:      []interface{}{accounts.HexToAddress(pairContractAddr), to.AmountIn},
	}
	transferTx, err := u.Web3Client.TransferERC20Token(ctx, transferTxParams)
	if err != nil {
		log.Warn().Interface("transferTx", transferTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(accounts.Hash(transferTx.Hash()))
	return err
}

func (u *UniswapClient) ApproveSpender(ctx context.Context, tokenAddress, spenderAddr string, amount *big.Int) (*types.Transaction, error) {
	approveTx, err := u.Web3Client.ERC20ApproveSpender(ctx, tokenAddress, spenderAddr, amount)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving spender")
		return approveTx, err
	}
	return approveTx, err
}

func (w *Web3Client) SignAndSendSmartContractTxPayload(ctx context.Context, scInfo *web3_actions.SendContractTxPayload) (*types.Transaction, error) {
	// TODO improve gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := w.GetSignedTxToCallFunctionWithArgs(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	err = w.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		log.Err(err).Msg("SignAndSendSmartContractTxPayload: failed to send signed tx")
		return nil, err
	}
	return signedTx, nil
}
