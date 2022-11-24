package ecdsa_signer

import (
	"encoding/hex"

	"github.com/gochain/gochain/v4/crypto"
	"github.com/gochain/web3"
)

type EcdsaSigner struct {
	*web3.Account
}

func CreateEcdsaSignerFromPk(pk string) (EcdsaSigner, error) {
	e := EcdsaSigner{}
	a, err := web3.ParsePrivateKey(pk)
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
