package artemis_trading_constants

import "github.com/zeus-fyi/gochain/web3/accounts"

const (
	Doge2ContractAddr    = "0xF2ec4a773ef90c58d98ea734c0eBDB538519b988"
	PepeContractAddr     = "0x6982508145454Ce325dDbE47a25d4ec3d2311933"
	TetherContractAddr   = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	BinanceCoinAddr      = "0xB8c77482e45F1F44dE1745F52C74426C631bDD52"
	UsdCoinAddr          = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	LidoSEthAddr         = "0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84"
	MaticTokenAddr       = "0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0"
	WETH9ContractAddress = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	FourCoinAddr         = "0x244b797d622D4DEe8b188b03546ACAAbD0Cf91A0"
	PoohCoinAddr         = "0xB69753c06BB5c366BE51E73bFc0cC2e3DC07E371"
	BobTokenAddr         = "0x7D8146cf21e8D7cbe46054e01588207b51198729"
	MongCoinAddr         = "0x1ce270557C1f68Cfb577b856766310Bf8B47FD9C"
	WojakTokenAddr       = "0x5026F006B85729a8b14553FAE6af249aD16c9aaB"
	GyoshiTokenAddr      = "0x1F17D72cBe65Df609315dF5c4f5F729eFbd00Ade"
	PPizzaTokenAddr      = "0xab306326bC72c2335Bd08F42cbec383691eF8446"
	WassieTokenAddr      = "0x2c95D751DA37A5C1d9c5a7Fd465c1d50F3d96160"
	CapyBaraTokenAddr    = "0xF03D5fC6E08dE6Ad886fCa34aBF9a59ef633b78a"
	FtmTokenAddr         = "0x4E15361FD6b4BB609Fa63C81A2be19d873717870"
	TurboTokenAddr       = "0xA35923162C49cF95e6BF26623385eb431ad920D3"
	SpongeTokenAddr      = "0x25722Cd432d02895d9BE45f5dEB60fc479c8781E"
	WBTCContractAddr     = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
	FraxTokenAddr        = "0x853d955aCEf822Db058eb8505911ED77F175b99e"
	LooksTokenAddr       = "0xf4d2888d29D722226FafA5d9B24F9164c092421E"
	BlurTokenAddr        = "0x5283D291DBCF85356A21bA090E6db59121208b44"
	ApeCoinTokenAddr     = "0x4d224452801ACEd8B2F0aebE155379bb5D594381"
	LinkTokenAddr        = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	// 8 decimal places
	HexTokenAddr                = "0x2b591e99afE9f32eAA6214f7B7629768c40Eeb39"
	DaiContractAddress          = "0x6b175474e89094c44da98b954eedeac495271d0f"
	Permit2SmartContractAddress = "0x000000000022D473030F116dDEE9F6B43aC78BA3"

	ZeroAddress       = "0x0000000000000000000000000000000000000000"
	Multicall3Address = "0xcA11bde05977b3631167028862bE2a173976CA11"

	MogToken = "0xaaee1a9723aadb7afa2810263653a34ba2c21c7a"
	BOBO     = "0xB90B2A35C65dBC466b04240097Ca756ad2005295"
)

var (
	Permit2SmartContractAddressAccount = accounts.HexToAddress(Permit2SmartContractAddress)
	Doge2ContractAddrAccount           = accounts.HexToAddress(Doge2ContractAddr)
	PepeContractAddrAccount            = accounts.HexToAddress(PepeContractAddr)
	WETH9ContractAddressAccount        = accounts.HexToAddress(WETH9ContractAddress)
	DaiContractAddressAccount          = accounts.HexToAddress(DaiContractAddress)
	Multicall3AddressAccount           = accounts.HexToAddress(Multicall3Address)
	LinkTokenAddressAccount            = accounts.HexToAddress(LinkTokenAddr)
	MogTokenAddressAccount             = accounts.HexToAddress(MogToken)
	BoboTokenAddressAccount            = accounts.HexToAddress(BOBO)
)

const (
	GoerliWETH9ContractAddress   = "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6"
	GoerliUniswapContractAddress = "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
	GoerliDaiContractAddress     = "0xdc31ee1784292379fbb2964b3b9c4124d8f89c60"
)

var (
	GoerliWETH9ContractAddressAccount   = accounts.HexToAddress(GoerliWETH9ContractAddress)
	GoerliUniswapContractAddressAccount = accounts.HexToAddress(GoerliUniswapContractAddress)
	GoerliDaiContractAddressAccount     = accounts.HexToAddress(GoerliDaiContractAddress)
)
