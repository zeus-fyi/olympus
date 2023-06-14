package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (s *Web3ClientTestSuite) TestPermit2() {
	expiration, _ := new(big.Int).SetString("946902158100", 10)
	sigDeadline, _ := new(big.Int).SetString("146902158100", 10)

	amount, _ := new(big.Int).SetString("100", 10)
	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				TokenPermissions: TokenPermissions{
					Token:  accounts.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"),
					Amount: amount,
				},
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress("0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"),
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}
	b := _hashPermitDetails(pp.PermitSingle.PermitDetails)
	hash := common.BytesToHash(b)
	exphash := common.HexToHash("0xc87aa0e9fdf4af6f31d56f7ed46715f6baba8e8f1ffdb494118f0f8b23f02c69")
	s.Equal(exphash, hash)
	eip := NewEIP712ForPermit2(chainID, accounts.HexToAddress("0x00001f78189bE22C3498cFF1B8e02272C3220000"))
	hashed := hashPermitSingle(pp.PermitSingle)
	s.Equal("0xa90c13eed97d34532a906c39ae1c798a831c8e26acd74c8e12008fed69aded02", hashed.String())
	hashed = eip.HashTypedData(hashed)
	s.Equal("0x6c6d214e1929d8e54266ff0eacff92afbc298d02316fe0cb4d9cdad1dea889b2", hashed.String())

}
