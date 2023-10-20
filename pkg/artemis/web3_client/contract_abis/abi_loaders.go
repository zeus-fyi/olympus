package artemis_oly_contract_abis

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	apps_hardhat "github.com/zeus-fyi/olympus/apps/olympus/hardhat"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var (
	ctx                    = context.Background()
	MultiCall3             = MustLoadMulticall3Abi()
	Erc20                  = MustLoadERC20Abi()
	Permit2                = MustLoadPermit2Abi()
	UniversalRouterDecoder = MustLoadUniversalRouterDecodingAbi()
	UniversalRouterNew     = MustLoadNewUniversalRouterAbi()
	UniversalRouterOld     = MustLoadNewUniversalRouterAbi()
	UniswapV2Router01      = MustLoadUniswapV2Router01ABI()
	UniswapV2Router02      = MustLoadUniswapV2Router02ABI()
	UniswapV3Router01      = MustLoadUniswapV3Swap1RouterAbi()
	UniswapV3Router02      = MustLoadUniswapV3Swap2RouterAbi()
)

func MustLoadMulticall3Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(Multicall3Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadPepeAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(PepeAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniswapV2PairAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapV2PairAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniswapV3Swap1RouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapV3RouterSwapV1ABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniswapV3Swap2RouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapV3RouterSwapV2ABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadQuoterV1Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(QuoterV1Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadTickLensAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(TickLensAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadQuoterV2Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(QuoterV2Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadPoolV3Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapPoolV3Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniswapV2Router02ABI() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapV2Router02ABI))
	if err != nil {
		panic(err)
	}
	return readAbi

}
func MustLoadUniswapV2Router01ABI() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniswapV2Router01ABI))
	if err != nil {
		panic(err)
	}
	return readAbi

}
func MustLoadOldUniversalRouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniversalRouterAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadNewUniversalRouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniversalRouterNewAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniversalRouterDecodingAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(UniversalRouterDecodingAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadPermit2Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(Permit2Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadERC20Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(Erc20Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}
func ForceDirToTestDirLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func LoadNewERC20AbiPayload() (web3_actions.SendContractTxPayload, error) {
	ForceDirToTestDirLocation()
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		FnIn:        "tmp.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, err
	}
	abiInput := m["ABI"]
	abf := &abi.ABI{}
	err = abf.UnmarshalJSON([]byte(abiInput.(string)))
	if err != nil {
		return web3_actions.SendContractTxPayload{}, err
	}
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, nil
}

func LoadERC20AbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "Token.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abiInput := m["abi"]
	b, err := json.Marshal(abiInput)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abf := &abi.ABI{}
	err = abf.UnmarshalJSON(b)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}

func LoadERC20DeployedByteCode() (string, error) {
	ForceDirToTestDirLocation()

	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		FnIn:        "Token.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return "", err
	}
	return m["deployedBytecode"].(string), nil
}

func LoadUniswapFactoryAbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "UniswapV2Factory.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abiInput := m["abi"]
	b, err := json.Marshal(abiInput)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abf := &abi.ABI{}
	err = abf.UnmarshalJSON(b)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}

func MustLoadSwapAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(SwapABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func LoadSwapAbiPayload(pairContractAddr string) (web3_actions.SendContractTxPayload, string, error) {
	abf := MustLoadSwapAbi()
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: pairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractFile:      "",
		ContractABI:       abf,
		MethodName:        "swap",
		Params:            []interface{}{},
	}
	return params, "", nil
}

func MustLoadRawdawgAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(RawdawgAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadRawdawgContractDeployPayload() (web3_actions.SendContractTxPayload, string) {
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      MustLoadRawdawgAbi(),
		Params:           []interface{}{},
	}
	return params, RawdawgByteCode
}

func LoadLocalRawdawgAbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	apps_hardhat.ForceDirToLocation()
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./artifacts/contracts/RawDawg.sol",
		FnIn:        "Rawdawg.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abiInput := m["abi"]
	b, err := json.Marshal(abiInput)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abf := &abi.ABI{}
	err = abf.UnmarshalJSON(b)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	//RawdawgAbi = abf
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}
