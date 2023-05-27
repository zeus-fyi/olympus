package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (w *Web3Client) ERC20ApproveSpender(ctx context.Context, scAddr, spenderAddr string, amount *big.Int) (*web3_types.Transaction, error) {
	w.Dial()
	defer w.Close()
	abiFile := LoadERC20Abi()
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: scAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       abiFile,
		Params:            []interface{}{spenderAddr, amount},
	}
	tx, err := w.ApproveSpenderERC20Token(ctx, payload)
	if err != nil {
		fmt.Println("error", err)
		return tx, err
	}
	return tx, err
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
		resp, err := w.GetStorageAt(ctx, scAddr, hexStr)
		if err != nil {
			fmt.Println("error", err)
			return -1, "", err
		}
		val := new(big.Int).SetBytes(resp)
		if val.String() == b.String() {
			return i, hexStr, nil
		}
	}
	return -1, "", errors.New("no slot found")
}

func (w *Web3Client) SetERC20BalanceAtSlotNumber(ctx context.Context, scAddr, userAddr string, slotNum int, value *big.Int) error {
	slotHex, err := getSlot(userAddr, new(big.Int).SetUint64(uint64(slotNum)))
	if err != nil {
		return err
	}
	newBalance := common.LeftPadBytes(value.Bytes(), 32)
	err = w.SetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
	if err != nil {
		return err
	}
	return nil
}

func (w *Web3Client) SetERC20BalanceBruteForce(ctx context.Context, scAddr, userAddr string, value *big.Int) error {
	for i := 0; i < 100; i++ {
		slotHex, err := getSlot(userAddr, new(big.Int).SetUint64(uint64(i)))
		if err != nil {
			return err
		}
		newBalance := common.LeftPadBytes(value.Bytes(), 32)
		err = w.SetStorageAt(ctx, scAddr, slotHex, common.BytesToHash(newBalance).Hex())
		if err != nil {
			continue
		}
		b, err := w.ReadERC20TokenBalance(ctx, scAddr, userAddr)
		if err != nil {
			return err
		}
		if b.String() == value.String() {
			return nil
		}
	}
	return errors.New("unable to overwrite balance")
}

func (w *Web3Client) MatchFrontRunTradeValues(tf *TradeExecutionFlowInBigInt) error {
	err := w.SetERC20BalanceBruteForce(ctx, tf.FrontRunTrade.AmountInAddr.String(), w.PublicKey(), tf.FrontRunTrade.AmountIn)
	if err != nil {
		return err
	}
	b, err := w.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), w.PublicKey())
	if err != nil {
		return err
	}
	if b.String() != tf.FrontRunTrade.AmountIn.String() {
		return errors.New("amount in not set correctly")
	}
	return nil
}
