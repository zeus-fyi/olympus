package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

type TokenFee struct {
	Token accounts.Address `json:"token"`
	Fee   *big.Int         `json:"fee"`
}

type TokenFeePath struct {
	TokenIn accounts.Address `json:"tokenIn"`
	Path    []TokenFee       `json:"path"`
}

func (tfp *TokenFeePath) GetFirstFee() *big.Int {
	if len(tfp.Path) == 0 {
		return nil
	}
	return tfp.Path[0].Fee
}

func (tfp *TokenFeePath) Encode() []byte {
	// Convert TokenIn into bytes
	tokenIn := tfp.TokenIn.Bytes()

	// Initialize a slice to hold the path bytes
	var pathBytes []byte

	// Iterate over the path to encode each TokenFee into bytes
	for _, tf := range tfp.Path {
		// Convert each TokenFee's token into bytes
		token := tf.Token.Bytes()

		// Convert each TokenFee's fee into a 3 bytes (6 hex characters)
		feeBytes := big.NewInt(tf.Fee.Int64()).Bytes()
		// If feeBytes is not 3 bytes long, pad it with leading zeros
		for len(feeBytes) < 3 {
			feeBytes = append([]byte{0}, feeBytes...)
		}

		// Append the fee and token bytes to the pathBytes
		pathBytes = append(pathBytes, feeBytes...)
		pathBytes = append(pathBytes, token...)
	}

	// Concatenate TokenIn and Path bytes
	return append(tokenIn, pathBytes...)
}

func (tfp *TokenFeePath) GetEndToken() accounts.Address {
	return tfp.Path[len(tfp.Path)-1].Token
}

func (tfp *TokenFeePath) GetPath() []accounts.Address {
	path := []accounts.Address{tfp.TokenIn}
	for _, p := range tfp.Path {
		path = append(path, p.Token)
	}
	return path
}

func (tfp *TokenFeePath) Reverse() {
	pathList := tfp.Path
	for i, j := 0, len(pathList)-1; i < j; i, j = i+1, j-1 {
		pathList[i], pathList[j] = pathList[j], pathList[i]
	}
	pathList[len(pathList)-1].Token, tfp.TokenIn = tfp.TokenIn, pathList[len(pathList)-1].Token
}
