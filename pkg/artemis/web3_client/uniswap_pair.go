package web3_client

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	getReserves       = "getReserves"
	pairAddressSuffix = "96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"
	ZeroEthAddress    = "0x0000000000000000000000000000000000000000"
)

type UniswapV2Pair struct {
	PairContractAddr     string           `json:"pairContractAddr"`
	Price0CumulativeLast *big.Int         `json:"price0CumulativeLast"`
	Price1CumulativeLast *big.Int         `json:"price1CumulativeLast"`
	KLast                *big.Int         `json:"kLast"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             *big.Int         `json:"reserve0"`
	Reserve1             *big.Int         `json:"reserve1"`
	BlockTimestampLast   *big.Int         `json:"blockTimestampLast"`
}

func (p *UniswapV2Pair) PairForV2FromAddresses(tokenA, tokenB accounts.Address) error {
	return p.PairForV2(tokenA.String(), tokenB.String())
}

func (p *UniswapV2Pair) PairForV2(tokenA, tokenB string) error {
	if tokenA == ZeroEthAddress {
		tokenA = WETH9ContractAddress
	}
	if tokenB == ZeroEthAddress {
		tokenB = WETH9ContractAddress
	}
	if tokenA == tokenB {
		return errors.New("identical addresses")
	}
	p.sortTokens(accounts.HexToAddress(tokenA), accounts.HexToAddress(tokenB))
	message := []byte{255}
	message = append(message, common.HexToAddress(UniswapV2FactoryAddress).Bytes()...)
	addrSum := p.Token0.Bytes()
	addrSum = append(addrSum, p.Token1.Bytes()...)
	message = append(message, crypto.Keccak256(addrSum)...)
	b, err := hex.DecodeString(pairAddressSuffix)
	if err != nil {
		return err
	}
	message = append(message, b...)
	hashed := crypto.Keccak256(message)
	addressBytes := big.NewInt(0).SetBytes(hashed)
	addressBytes = addressBytes.Abs(addressBytes)
	p.PairContractAddr = common.BytesToAddress(addressBytes.Bytes()).String()
	return nil
}

func (u *UniswapClient) V2PairToPrices(ctx context.Context, pairAddr []accounts.Address) (UniswapV2Pair, error) {
	p := UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return p, err
		}
		err = u.GetPairContractPrices(ctx, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return p, err
		}
		return p, err
	}
	return UniswapV2Pair{}, errors.New("pair address length is not 2, multi-hops not implemented yet")
}

func (p *UniswapV2Pair) GetQuoteUsingTokenAddr(addr string, amount *big.Int) (*big.Int, error) {
	if addr == "0x0000000000000000000000000000000000000000" {
		addr = WETH9ContractAddress
	}
	if p.Token0 == accounts.HexToAddress(addr) {
		return p.GetQuoteToken0BuyToken1(amount)
	}
	if p.Token1 == accounts.HexToAddress(addr) {
		return p.GetQuoteToken1BuyToken0(amount)
	}
	return nil, errors.New("GetQuoteUsingTokenAddr: token not found")
}

func (p *UniswapV2Pair) GetOppositeToken(addr string) accounts.Address {
	if addr == "0x0000000000000000000000000000000000000000" {
		addr = WETH9ContractAddress
	}
	if p.Token0 == accounts.HexToAddress(addr) {
		return p.Token1
	}
	if p.Token1 == accounts.HexToAddress(addr) {
		return p.Token0
	}
	log.Warn().Msgf("GetOppositeToken: token not found: %s", addr)
	return accounts.Address{}
}

func (p *UniswapV2Pair) GetTokenNumber(addr accounts.Address) int {
	if addr.String() == "0x0000000000000000000000000000000000000000" {
		if p.Token0.String() == WETH9ContractAddress {
			return 0
		}
		if p.Token1.String() == WETH9ContractAddress {
			return 1
		}
	}
	if p.Token0 == addr {
		return 0
	}
	if p.Token1 == addr {
		return 1
	}
	log.Warn().Msgf("GetTokenNumber: token not found: %s", addr)
	return -1
}

func (u *UniswapClient) GetPairContractPrices(ctx context.Context, p *UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.PairAbi,
	}
	scInfo.MethodName = getReserves
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return err
	}
	if len(resp) <= 2 {
		return err
	}
	reserve0, err := ParseBigInt(resp[0])
	if err != nil {
		return err
	}
	p.Reserve0 = reserve0
	reserve1, err := ParseBigInt(resp[1])
	if err != nil {
		return err
	}
	p.Reserve1 = reserve1
	blockTimestampLast, err := ParseBigInt(resp[2])
	if err != nil {
		return err
	}
	p.BlockTimestampLast = blockTimestampLast
	return nil
}

/*
	price0, err := u.SingleReadMethodBigInt(ctx, "price0CumulativeLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Price0CumulativeLast = price0
	price1, err := u.SingleReadMethodBigInt(ctx, "price1CumulativeLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Price1CumulativeLast = price1
	kLast, err := u.SingleReadMethodBigInt(ctx, "kLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.KLast = kLast
	token0, err := u.SingleReadMethodAddr(ctx, "token0", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Token0 = token0
	token1, err := u.SingleReadMethodAddr(ctx, "token1", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
*/
