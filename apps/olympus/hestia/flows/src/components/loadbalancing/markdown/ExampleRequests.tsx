import {Groups} from "../../../redux/loadbalancing/loadbalancing.types";


export function findKeyWithPrefix(groups: Groups): string  {
    let fallbackKey = ''

    for (const key in groups) {
        switch (true) {
            case key.startsWith('ethereum') || key.startsWith('celo') || key.startsWith('polygon')
            || key.startsWith('zk') || key.startsWith('base') || key.startsWith('opt') || key.startsWith('arb') || key.startsWith('bsc') || key.startsWith('bnb'):
                if (key.startsWith(key)) return key;
                break;
            case key.startsWith('btc'):
                if (key.startsWith(key)) return key;
                break;
            case key.startsWith('near'):
                if (key.startsWith(key)) return key;
                break;
            case key.startsWith('avalanche'):
                if (key.startsWith(key)) return key;
                break;
        }
        if (key !== '-all' && key !== 'unused') {
            fallbackKey = key;
        }
    }
    return fallbackKey;
}

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

