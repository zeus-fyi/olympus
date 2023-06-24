package web3_client

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	apps_hardhat "github.com/zeus-fyi/olympus/apps/olympus/hardhat"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func MustLoadUniswapV3RouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV3RouterABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadQuoterV1Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.QuoterV1Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadTickLensAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.TickLensAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadQuoterV2Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.QuoterV2Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadPoolV3Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapPoolV3Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniswapV2RouterABI() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2RouterABI))
	if err != nil {
		panic(err)
	}
	return readAbi

}
func MustLoadUniversalRouterAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniversalRouterAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadUniversalRouterDecodingAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniversalRouterDecodingAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadPermit2Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.Permit2Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func MustLoadERC20Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.Erc20Abi))
	if err != nil {
		panic(err)
	}
	return readAbi
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
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
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
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.SwapABI))
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
		MethodName:        swap,
		Params:            []interface{}{},
	}
	return params, "", nil
}

func MustLoadRawdawgAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.RawdawgAbi))
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
	return params, artemis_oly_contract_abis.RawdawgByteCode
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
	RawdawgAbi = abf
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}
