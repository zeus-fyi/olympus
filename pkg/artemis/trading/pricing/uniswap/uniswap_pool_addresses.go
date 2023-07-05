package uniswap_pricing

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/utils"
)

func (u *UniswapPools) getPoolAddress(pairAddr []accounts.Address) error {
	err := u.V2Pair.PairForV2(pairAddr[0].String(), pairAddr[1].String())
	if err != nil {
		log.Err(err).Msg("GetPoolAddress: PairForV2")
		return err
	}
	err = u.PairForV3FromAddresses(pairAddr[0], pairAddr[1])
	if err != nil {
		return err
	}
	return nil
}

type UniswapPools struct {
	V2Pair  UniswapV2Pools
	V3Pairs UniswapV3Pools
}

func NewUniswapPools(pair []accounts.Address) (*UniswapPools, error) {
	if len(pair) != 2 {
		log.Err(errors.New("pair address length is not 2, multi-hops not implemented yet")).Msg("GetPoolAddress")
		return nil, errors.New("pair address length is not 2, multi-hops not implemented yet")
	}
	pools := UniswapPools{
		V2Pair: UniswapV2Pools{
			UniswapV2Pair: &UniswapV2Pair{},
		},
		V3Pairs: UniswapV3Pools{
			LowFee: &UniswapV3Pair{
				Fee: constants.FeeLow,
			},
			MediumFee: &UniswapV3Pair{
				Fee: constants.FeeMedium,
			},
			HighFee: &UniswapV3Pair{
				Fee: constants.FeeHigh,
			},
		},
	}
	err := pools.getPoolAddress(pair)
	if err != nil {
		return nil, err
	}
	return &pools, nil
}

type UniswapV2Pools struct {
	*UniswapV2Pair
}

type UniswapV3Pools struct {
	LowFee    *UniswapV3Pair
	MediumFee *UniswapV3Pair
	HighFee   *UniswapV3Pair
}

func (u *UniswapPools) PairForV3FromAddresses(token0, token1 accounts.Address) error {
	tokenA := uniswap_core_entities.NewToken(1, token0, 0, "", "")
	tokenB := uniswap_core_entities.NewToken(1, token1, 0, "", "")
	factoryAddress := artemis_trading_constants.UniswapV3FactoryAddressAccount
	pa, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeLow, "")
	if err != nil {
		log.Err(err).Msg("UniswapV3Pair: PairForV2")
		return err
	}
	u.V3Pairs.LowFee.PoolAddress = pa.String()
	pa, err = utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	if err != nil {
		log.Err(err).Msg("UniswapV3Pair: PairForV2")
		return err
	}
	u.V3Pairs.MediumFee.PoolAddress = pa.String()
	pa, err = utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeHigh, "")
	if err != nil {
		log.Err(err).Msg("UniswapV3Pair: PairForV2")
		return err
	}
	u.V3Pairs.HighFee.PoolAddress = pa.String()
	return err
}

func (p *UniswapV2Pair) PairForV2FromAddresses(tokenA, tokenB accounts.Address) error {
	return p.PairForV2(tokenA.String(), tokenB.String())
}

func (p *UniswapV2Pair) PairForV2(tokenA, tokenB string) error {
	if tokenA == ZeroEthAddress {
		tokenA = WETH.String()
	}
	if tokenB == ZeroEthAddress {
		tokenB = WETH.String()
	}
	if tokenA == tokenB {
		return errors.New("identical addresses")
	}
	p.sortTokens(accounts.HexToAddress(tokenA), accounts.HexToAddress(tokenB))
	message := []byte{255}
	message = append(message, common.HexToAddress(UniswapV2FactoryAddress.String()).Bytes()...)
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
