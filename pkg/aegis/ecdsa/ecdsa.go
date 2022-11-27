package ecdsa_signer

import (
	"encoding/hex"

	"github.com/gochain/gochain/v4/crypto"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type EcdsaSigner struct {
	*accounts.Account
}

func CreateEcdsaSignerFromPk(pk string) (EcdsaSigner, error) {
	e := EcdsaSigner{}
	a, err := accounts.ParsePrivateKey(pk)
	if err != nil {
		return e, err
	}
	e.Account = a
	return e, err
}

func NewEcdsaPkHexString() string {
	pk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	pkHexString := "0x" + hex.EncodeToString(crypto.FromECDSA(pk))
	return pkHexString
}
