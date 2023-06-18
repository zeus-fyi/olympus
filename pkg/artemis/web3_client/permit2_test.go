package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

const (
	maxUINT = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
)

func (s *Web3ClientTestSuite) TestCopyPermitTest() {
	expiration, _ := new(big.Int).SetString("3000000000000", 10)
	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				Token:      accounts.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
				Amount:     new(big.Int).SetUint64(1000000000),
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress("0xe808c1cfeebb6cb36b537b82fa7c9eef31415a05"),
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	permitAddress := "0x4a873bdd49f7f9cc0a5458416a12973fab208f8d"
	err := pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress(permitAddress), "Permit2")
	s.Require().Nil(err)
	s.Require().NotNil(pp.Signature)

	hashed := hashPermitSingle(pp.PermitSingle)
	eip := NewEIP712(chainID, accounts.HexToAddress(permitAddress), "Permit2")
	hashed = eip.HashTypedData(hashed)

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress(permitAddress), "Permit2")
	s.Require().Nil(err)

	verified, err := s.LocalHardhatMainnetUser.VerifySignature(s.LocalHardhatMainnetUser.Address(), hashed.Bytes(), pp.Signature)
	s.Require().Nil(err)
	s.Require().True(verified)

	jsSig := "1a622a5fb555e46f58b11ace6176bfc6d1f8ac4be3711612e5f89027de9aae96490d65fc3dce716c08cef58f1d78856fa0a50d13512cd207206d7aca11017ed100"
	s.Equal(jsSig, common.Bytes2Hex(pp.Signature))
}

func (s *Web3ClientTestSuite) TestPermit2TransferSubmission() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 17461010)
	s.Require().Nil(err)
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	wethAddress := accounts.HexToAddress(WETH9ContractAddress)
	permit2Address := accounts.HexToAddress(Permit2SmartContractAddress)

	fmt.Println(permit2Address.String())
	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, wethAddress.String(), s.LocalHardhatMainnetUser.PublicKey(), EtherMultiple(10000000))
	s.Require().Nil(err)

	bal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, wethAddress.String(), s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Require().NotNil(bal)
	fmt.Println(bal.String())
	fmt.Println(s.HostedHardhatMainnetUser.PublicKey())
	bal, err = s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
	s.Require().Nil(err)
	s.Require().NotNil(bal)
	fmt.Println(bal.String())
	nbal := hexutil.Big{}
	bigInt := nbal.ToInt()
	bigInt.Set(EtherMultiple(10000000))
	nbal = hexutil.Big(*bigInt)
	err = s.LocalHardhatMainnetUser.SetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nbal)
	s.Require().Nil(err)
	bal, err = s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
	s.Require().Nil(err)
	s.Require().NotNil(bal)
	max, _ := new(big.Int).SetString(maxUINT, 10)
	expiration, _ := new(big.Int).SetString("1785444080", 10)
	sigDeadline, _ := new(big.Int).SetString("1785444080", 10)

	tx, err := uni.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, max)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)
	bal, err = s.LocalHardhatMainnetUser.ReadERC20Allowance(ctx, wethAddress.String(), s.LocalHardhatMainnetUser.PublicKey(), Permit2SmartContractAddress)
	s.Assert().NoError(err)
	s.Assert().NotNil(bal)
	s.Require().Equal(max, bal)

	name, err := s.LocalHardhatMainnetUser.ReadERC20TokenName(ctx, wethAddress.String())
	s.Assert().NoError(err)
	s.Assert().Equal("Wrapped Ether", name)

	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				Token:      wethAddress,
				Amount:     EtherMultiple(1),
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     urAddr,
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, permit2Address, "Permit2")
	s.Assert().NoError(err)
	s.Assert().NotNil(pp.Signature)
	pp.Signature[64] += 27

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: Permit2SmartContractAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPermit2Abi(),
		MethodName:        permit0,
		Params:            []interface{}{s.LocalHardhatMainnetUser.Account.Address().String(), pp.PermitSingle, pp.Signature},
	}
	tx, err = s.LocalHardhatMainnetUser.SignAndSendSmartContractTxPayload(ctx, scInfo)
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)
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
				Token:      accounts.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"),
				Amount:     amount,
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

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress("0xCe71065D4017F316EC606Fe4422e11eB2c47c246"), "Permit2")
	s.NoError(err)
	s.Equal(sig, pp.Signature)
}
