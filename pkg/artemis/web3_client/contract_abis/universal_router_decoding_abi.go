package artemis_oly_contract_abis

const UniversalRouterDecodingAbi = `
[{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountIn",
      "type": "uint256"
    },
    {
      "name": "amountOutMin",
      "type": "uint256"
    },
    {
      "name": "path",
      "type": "address[]"
    },
    {
      "name": "payerIsUser",
      "type": "bool"
    }
  ],
  "name": "V3_SWAP_EXACT_IN",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountOut",
      "type": "uint256"
    },
    {
      "name": "amountInMax",
      "type": "uint256"
    },
    {
      "name": "path",
      "type": "address[]"
    },
    {
      "name": "payerIsUser",
      "type": "bool"
    }
  ],
  "name": "V3_SWAP_EXACT_OUT",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountIn",
      "type": "uint256"
    },
    {
      "name": "amountOutMin",
      "type": "uint256"
    },
    {
      "name": "path",
      "type": "address[]"
    },
    {
      "name": "payerIsSender",
      "type": "bool"
    }
  ],
  "name": "V2_SWAP_EXACT_IN",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountOut",
      "type": "uint256"
    },
    {
      "name": "amountInMax",
      "type": "uint256"
    },
    {
      "name": "path",
      "type": "address[]"
    },
    {
      "name": "payerIsSender",
      "type": "bool"
    }
  ],
  "name": "V2_SWAP_EXACT_OUT",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "token",
      "type": "address"
    },
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountMin",
      "type": "uint256"
    }
  ],
  "name": "SWEEP",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "token",
      "type": "address"
    },
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amount",
      "type": "uint256"
    }
  ],
  "name": "PERMIT2_TRANSFER_FROM",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountMin",
      "type": "uint256"
    }
  ],
  "name": "WRAP_ETH",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "amountMin",
      "type": "uint256"
    }
  ],
  "name": "UNWRAP_WETH",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "token",
      "type": "address"
    },
    {
      "name": "recipient",
      "type": "address"
    },
    {
      "name": "value",
      "type": "uint256"
    }
  ],
  "name": "TRANSFER",
  "type": "function"
}]`
