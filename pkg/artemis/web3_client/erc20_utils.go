package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

var Erc20Abi = artemis_oly_contract_abis.MustLoadERC20Abi()

func (w *Web3Client) ERC20ApproveSpender(ctx context.Context, scAddr, spenderAddr string, amount *big.Int) (*types.Transaction, error) {
	w.Dial()
	defer w.Close()

	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: scAddr,
		MethodName:        "approve",
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs:   web3_actions.TransferArgs{},
			GasPriceLimits: web3_actions.GasPriceLimits{},
		},
		ContractABI: Erc20Abi,
		Params:      []interface{}{accounts.HexToAddress(spenderAddr), amount},
	}
	signedTx, err := w.CallFunctionWithArgs(ctx, &payload)
	if err != nil {
		return nil, err
	}
	return signedTx, err
}

func (w *Web3Client) ERC20ApproveSpenderSignedTx(ctx context.Context, scAddr, spenderAddr string, amount *big.Int) (*types.Transaction, error) {
	w.Dial()
	defer w.Close()

	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: scAddr,
		MethodName:        "approve",
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs:   web3_actions.TransferArgs{},
			GasPriceLimits: web3_actions.GasPriceLimits{},
		},
		ContractABI: Erc20Abi,
		Params:      []interface{}{accounts.HexToAddress(spenderAddr), amount},
	}
	signedTx, err := w.GetSignedTxToCallFunctionWithArgs(ctx, &payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ERC20ApproveSpenderSignedTx: GetSignedTxToCallFunctionWithArgs")
		return nil, err
	}
	return signedTx, err
}

func (w *Web3Client) FindSlotFromUserWithBalance(ctx context.Context, scAddr, userAddr string) (int, string, error) {
	b, rerr := w.ReadERC20TokenBalance(ctx, scAddr, userAddr)
	if rerr != nil {
		return -1, "", rerr
	}
	if b.String() == "0" {
		return -1, "", errors.New("no balance")
	}
	for i := 0; i < 100; i++ {
		slotNum := new(big.Int).SetUint64(uint64(i))
		hexStr, err := GetSlot(userAddr, slotNum)
		if err != nil {
			return -1, "", err
		}
		resp, err := w.HardHatGetStorageAt(ctx, scAddr, hexStr)
		if err != nil {
			fmt.Println("error", err)
			return -1, "", err
		}
		val := new(big.Int).SetBytes(resp)
		if val.String() == b.String() {
			return i, hexStr, nil
		}
		time.Sleep(25 * time.Millisecond)
	}
	return -1, "", errors.New("no slot found")
}

func (w *Web3Client) SetERC20BalanceAtSlotNumber(ctx context.Context, scAddr, userAddr string, slotNum int, value *big.Int) error {
	slotHex, err := GetSlot(userAddr, new(big.Int).SetUint64(uint64(slotNum)))
	if err != nil {
		return err
	}
	newBalance := common.LeftPadBytes(value.Bytes(), 32)
	err = w.HardhatSetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
	if err != nil {
		return err
	}
	return nil
}

func (w *Web3Client) SetERC20BalanceBruteForce(ctx context.Context, scAddr, userAddr string, value *big.Int) error {
	// TODO assumes only mainnet for now
	protocolNetworkID := 1
	slotNum, serr := artemis_mev_models.SelectERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
		Address:           scAddr,
		ProtocolNetworkID: protocolNetworkID,
	})
	if serr != nil {
		return serr
	}

	if slotNum > -1 {
		slotHex, err := GetSlot(userAddr, new(big.Int).SetUint64(uint64(slotNum)))
		if err != nil {
			return err
		}
		newBalance := common.LeftPadBytes(value.Bytes(), 32)
		err = w.HardhatSetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
		return err
	}

	for i := 0; i < 20; i++ {
		slotHex, err := GetSlot(userAddr, new(big.Int).SetUint64(uint64(i)))
		if err != nil {
			return err
		}
		newBalance := common.LeftPadBytes(value.Bytes(), 32)
		err = w.HardhatSetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
		if err != nil {
			continue
		}
		b, err := w.ReadERC20TokenBalance(ctx, scAddr, userAddr)
		if err != nil {
			return err
		}
		if b.String() == value.String() {
			err = artemis_mev_models.UpdateERC20TokenBalanceOfSlotInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
				Address:           scAddr,
				ProtocolNetworkID: 1,
				BalanceOfSlotNum:  i,
			})
			if err != nil {
				log.Err(err)
				return err
			}
			return nil
		}
		time.Sleep(25 * time.Millisecond)
	}

	err := artemis_mev_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
		Address:           scAddr,
		ProtocolNetworkID: 1,
		BalanceOfSlotNum:  -1,
	})
	if err != nil {
		log.Err(err).Msg("error inserting token info")
	}
	return errors.New("unable to overwrite balance")
}

func (w *Web3Client) MatchFrontRunTradeValues(tf *TradeExecutionFlow) error {
	pubkey := w.PublicKey()
	err := w.SetERC20BalanceBruteForce(ctx, tf.FrontRunTrade.AmountInAddr.String(), pubkey, tf.FrontRunTrade.AmountIn)
	if err != nil {
		return err
	}
	b, err := w.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), pubkey)
	if err != nil {
		return err
	}
	if b.String() != tf.FrontRunTrade.AmountIn.String() {
		return errors.New("amount in not set correctly")
	}
	return nil
}

func (w *Web3Client) GetMainnetBalanceWETH(address string) (*big.Int, error) {
	b, err := w.ReadERC20TokenBalance(ctx, artemis_trading_constants.WETH9ContractAddress, address)
	if err != nil {
		return nil, err
	}
	return b, nil
}

const (
	ZeusTestSessionLockHeaderValue = "Zeus-Test"
	IrisAnvil                      = "https://iris.zeus.fyi/v1beta/internal/"
)

func (w *Web3Client) GetMainnetBalanceDiffWETH(address string, blockNumber int) (*big.Int, error) {
	w.IsAnvilNode = true
	w.NodeURL = IrisAnvil

	if w.GetSessionLockHeader() == "" {
		w.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	}
	if w.Headers == nil || w.Headers["Authorization"] == "" {
		return nil, errors.New("no bearer token")
	}
	w.Dial()
	defer w.Close()
	err := w.HardHatResetNetwork(ctx, blockNumber-1)
	if err != nil {
		return nil, err
	}
	startBal, err := w.ReadERC20TokenBalance(context.Background(), artemis_trading_constants.WETH9ContractAddress, address)
	if err != nil {
		return nil, err
	}
	err = w.HardHatResetNetwork(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	endBal, err := w.ReadERC20TokenBalance(context.Background(), artemis_trading_constants.WETH9ContractAddress, address)
	if err != nil {
		return nil, err
	}
	return artemis_eth_units.SubBigInt(endBal, startBal), nil
}

// only for debugging, not for usage that's not done manually

func (w *Web3Client) ResetNetworkLocalToExtIrisTest(blockNumber int) error {
	w.IsAnvilNode = true
	w.NodeURL = IrisAnvil
	w.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)

	if w.Headers == nil || w.Headers["Authorization"] == "" {
		return errors.New("no bearer token")
	}
	w.Dial()
	defer w.Close()
	err := w.HardHatResetNetwork(ctx, blockNumber)
	if err != nil {
		return err
	}
	return nil
}

//func (w *Web3Client) MatchFrontRunTradeValues(tf *TradeExecutionFlow) error {
//	pubkey := w.PublicKey()
//	slotHex, err := GetSlot(pubkey, new(big.Int).SetUint64(uint64(0)))
//	if err != nil {
//		return err
//	}
//	inAddr := tf.FrontRunTrade.AmountInAddr.String()
//	value := tf.FrontRunTrade.AmountIn
//	newBalance := common.LeftPadBytes(value.Bytes(), 32)
//	bc, err := artemis_oly_contract_abis.LoadERC20DeployedByteCode()
//	if err != nil {
//		return err
//	}
//	err = w.SetCodeOverride(ctx, inAddr, bc)
//	if err != nil {
//		return err
//	}
//	err = w.HardhatSetStorageAt(ctx, inAddr, slotHex, common.BytesToHash(newBalance).Hex())
//	if err != nil {
//		return err
//	}
//	b, err := w.ReadERC20TokenBalance(ctx, inAddr, pubkey)
//	if err != nil {
//		return err
//	}
//	if b.String() != tf.FrontRunTrade.AmountIn.String() {
//		return errors.New("amount in not set correctly")
//	}
//	return nil
//}
