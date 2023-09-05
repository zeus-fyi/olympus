package iris_catalog_procedures

const (
	EthMaxBlockAggReduce  = "eth_maxBlockAggReduce"
	AvaxMaxBlockAggReduce = "avax_maxBlockAggReduce"
	NearMaxBlockAggReduce = "near_maxBlockAggReduce"
	BtcMaxBlockAggReduce  = "btc_maxBlockAggReduce"
)

func ProcedureStageOnePayload(procName string) string {
	switch procName {
	case EthMaxBlockAggReduce:
		return EthGetBlockNumberPayload
	case AvaxMaxBlockAggReduce:
		return ""
	case NearMaxBlockAggReduce:
		return NearGetBlockNumberPayload
	case BtcMaxBlockAggReduce:
		return BtcGetBlockNumberPayload
	default:
		return ""
	}
}

const (
	EthGetBlockNumberPayload = `
	{
	  "method": "eth_blockNumber",
	  "params": [],
	  "id": 1,
	  "jsonrpc": "2.0"
	}`
	NearGetBlockNumberPayload = `
	{
		"jsonrpc": "2.0",
		"method": "status",
		"params": [],
		"id": 1
	}`
	BtcGetBlockNumberPayload = `{ "method": "getblockcount" }`
)
