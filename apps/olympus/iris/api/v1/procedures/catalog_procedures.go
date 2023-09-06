package iris_catalog_procedures

import "github.com/labstack/echo/v4"

const (
	EthMaxBlockAggReduce               = "eth_maxBlockAggReduce"
	NearMaxBlockAggReduce              = "near_maxBlockAggReduce"
	BtcMaxBlockAggReduce               = "btc_maxBlockAggReduce"
	AvaxContractChainMaxBlockAggReduce = "avax_maxBlockAggReduce"
	AvaxPlatformChainMaxBlockAggReduce = "avax_platformMaxHeightAggReduce"
)

func ProcedureStageOnePayload(procName string) echo.Map {
	switch procName {
	case EthMaxBlockAggReduce:
		return EthGetBlockNumberPayload
	case AvaxContractChainMaxBlockAggReduce:
		return AvaxContractChainGetBlockNumberPayload
	case AvaxPlatformChainMaxBlockAggReduce:
		return AvaxPlatformChainGetBlockNumberPayload
	case NearMaxBlockAggReduce:
		return NearGetBlockNumberPayload
	case BtcMaxBlockAggReduce:
		return BtcGetBlockNumberPayload
	default:
		return echo.Map{}
	}
}

func IsEmbeddedProcedure(procName string) bool {
	switch procName {
	case EthMaxBlockAggReduce:
		return true
	case AvaxContractChainMaxBlockAggReduce, AvaxPlatformChainMaxBlockAggReduce:
		return true
	case NearMaxBlockAggReduce:
		return true
	case BtcMaxBlockAggReduce:
		return true
	default:
		return false
	}
}

var (
	EthGetBlockNumberPayload = echo.Map{
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
		"jsonrpc": "2.0",
	}
	NearGetBlockNumberPayload = echo.Map{
		"jsonrpc": "2.0",
		"method":  "status",
		"params":  []interface{}{},
		"id":      1,
	}
	BtcGetBlockNumberPayload = echo.Map{
		"method": "getblockcount",
	}
	AvaxPlatformChainGetBlockNumberPayload = echo.Map{
		"jsonrpc": "2.0",
		"method":  "platform.getHeight",
	}
	AvaxContractChainGetBlockNumberPayload = echo.Map{
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
		"jsonrpc": "2.0",
	}
)
