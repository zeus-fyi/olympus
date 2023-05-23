package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/crypto"
)

func getSlot(userAddress string, slot *big.Int) (string, error) {
	// compute keccak256 hash
	addr := common.HexToAddress(userAddress)
	hash := crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(slot.Bytes(), 32),
	)

	// return hex string of the hash
	return hash.Hex(), nil
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
