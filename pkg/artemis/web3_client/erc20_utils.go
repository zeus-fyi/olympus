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
			GasPriceLimits: web3_actions.GasPriceLimits{
				//GasLimit:  200000,
				//GasPrice:  GweiMultiple(50),
				//GasTipCap: GweiMultiple(2),
				//GasFeeCap: GweiMultiple(100),
			},
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
		hexStr, err := getSlot(userAddr, slotNum)
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
	slotHex, err := getSlot(userAddr, new(big.Int).SetUint64(uint64(slotNum)))
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
		slotHex, err := getSlot(userAddr, new(big.Int).SetUint64(uint64(slotNum)))
		if err != nil {
			return err
		}
		newBalance := common.LeftPadBytes(value.Bytes(), 32)
		err = w.HardhatSetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
		return err
	}

	for i := 0; i < 100; i++ {
		slotHex, err := getSlot(userAddr, new(big.Int).SetUint64(uint64(i)))
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
