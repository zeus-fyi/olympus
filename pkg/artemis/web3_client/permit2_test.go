package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *Web3ClientTestSuite) TestPermit2TransferSubmission() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	wethAddress := accounts.HexToAddress(WETH9ContractAddress)
	spender := accounts.HexToAddress(UniswapUniversalRouterAddress)

	err := s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, wethAddress.String(), s.LocalMainnetWeb3User.PublicKey(), TenThousandEther)
	s.Require().Nil(err)

	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	tx, err := uni.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, TenThousandEther)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	expiration, _ := new(big.Int).SetString("16785444080", 10)
	sigDeadline, _ := new(big.Int).SetString("16759542080", 10)
	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				TokenPermissions: TokenPermissions{
					Token:  wethAddress,
					Amount: EtherMultiple(1),
				},
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     spender,
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, spender)
	s.Assert().NoError(err)
	s.Assert().NotNil(pp.Signature)

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: Permit2SmartContractAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPermit2Abi(),
		MethodName:        permit0,
		Params:            []interface{}{s.LocalHardhatMainnetUser.Account.Address(), pp.PermitDetails, pp.Signature},
	}
	tx, err = s.LocalHardhatMainnetUser.SignAndSendSmartContractTxPayload(ctx, scInfo)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	/*
		Transfer approved tokens from one address to another
		Requires the from address to have approved at least the desired amount of tokens to msg.sender.
		from := ""
		to := ""
		amount := ""
		token := ""
	*/

	//scInfo = &web3_actions.SendContractTxPayload{
	//	SmartContractAddr: Permit2SmartContractAddress,
	//	SendEtherPayload:  web3_actions.SendEtherPayload{},
	//	ContractABI:       MustLoadPermit2Abi(),
	//	MethodName:        transferFrom0,
	//	Params:            []interface{}{s.LocalHardhatMainnetUser.Account.Address(), spender, pp.Amount, pp.Token},
	//}
	//tx, err = s.LocalHardhatMainnetUser.SignAndSendSmartContractTxPayload(ctx, scInfo)
	//s.Assert().NoError(err)
	//s.Assert().NotNil(tx)
}

func (s *Web3ClientTestSuite) TestPermit2Transfer() {
	sigDeadline, _ := new(big.Int).SetString("146902158100", 10)
	amount, _ := new(big.Int).SetString("100", 10)
	pt := PermitTransferFrom{
		TokenPermissions: TokenPermissions{
			Token:  accounts.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"),
			Amount: amount,
		},
		Nonce:       new(big.Int).SetUint64(0),
		SigDeadline: sigDeadline,
	}
	hash := _hashTokenPermissions(pt.TokenPermissions)
	s.Equal("73dffa388f7cfcea85654f48d7cd2ff5daf542e0b51bba732287bdd89e73b35c", common.Bytes2Hex(hash[:]))

	hashVal := hashPermitTransferFrom(pt, s.LocalHardhatMainnetUser.Address())
	s.Equal("0x9b9bc3959c07ca67947b15a7d6e7fcab56c8c17a5755d7852f6081a8917efb5d", hashVal.String())
}

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
	eip := NewEIP712ForPermit2(chainID, accounts.HexToAddress("0xCe71065D4017F316EC606Fe4422e11eB2c47c246"))
	val := eip.DomainSeparator()
	s.Equal("0xb7319fe24f5e0c062ca659214a8812519139f17ade16f660e6a77e2f558d6e1a", val.String())
	hashed := hashPermitSingle(pp.PermitSingle)
	s.Equal("0xa90c13eed97d34532a906c39ae1c798a831c8e26acd74c8e12008fed69aded02", hashed.String())
	hashed = eip.HashTypedData(hashed)
	s.Equal("0x6a4964621b8c850feebefb04dd997d9d109a807ec26f7fdc26282c3b2f0e2c74", hashed.String())
	sig, err := s.LocalHardhatMainnetUser.Sign(hashed.Bytes())
	s.NoError(err)
	verified, err := s.LocalHardhatMainnetUser.VerifySignature(s.LocalHardhatMainnetUser.Address(), hashed.Bytes(), sig)
	s.NoError(err)
	s.True(verified)

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress("0xCe71065D4017F316EC606Fe4422e11eB2c47c246"))
	s.NoError(err)
	s.Equal(sig, pp.Signature)
}
