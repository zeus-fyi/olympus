

export const ethMaxBlockAggReduceExample = `{
    "jsonrpc": "2.0",
    "procedure": "eth_maxBlockAggReduce",
    "method": "eth_getBlockByNumber",
    "params": ["latest", true],
    "id": 1
}`

export const nearMaxBlockAggReduceExample = `{
    "jsonrpc": "2.0",
    "procedure": "near_maxBlockAggReduce",
    "method": "block",
    "params": {
        "finality": "final"
    },
    "id": 1
}`

export const avaxMaxBlockAggReduceExample = `{
    "jsonrpc": "2.0",
    "procedure": "avax_maxBlockAggReduce",
    "method": "eth_gasPrice",
    "params": [],
    "id": 1
}`

export const avaxPlatformMaxBlockAggReduceExample = `{
    "jsonrpc": "2.0",
    "procedure": "avax_platformMaxHeightAggReduce",
    "method": "platform.getCurrentSupply",
    "params": [],
    "id": 1
}`

export const btcMaxBlockAggReduceExample = `{
    "jsonrpc": "2.0",
    "procedure": "btc_maxBlockAggReduce",
    "method": "getdifficulty",
    "id": 1
}`

