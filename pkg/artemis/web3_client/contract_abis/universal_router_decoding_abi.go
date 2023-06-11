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
      "type": "bytes"
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
      "type": "bytes"
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
      "name": "value",
      "type": "uint256"
    },
    {
      "name": "data",
      "type": "bytes"
    }
  ],
  "name": "SUDOSWAP",
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
      "name": "amount",
      "type": "uint160"
    },
    {
      "name": "expiration",
      "type": "uint48"
    },
    {
      "name": "nonce",
      "type": "uint48"
    },
    {
      "name": "spender",
      "type": "address"
    },
    {
      "name": "sigDeadline",
      "type": "uint256"
    },
    {
      "name": "signature",
      "type": "bytes"
    }
  ],
  "name": "PERMIT2_PERMIT",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "permitBatch",
      "type": "tuple",
      "components": [
        {
          "name": "details",
          "type": "tuple[]",
          "components": [
            {
              "name": "token",
              "type": "address"
            },
            {
              "name": "amount",
              "type": "uint160"
            },
            {
              "name": "expiration",
              "type": "uint48"
            },
            {
              "name": "nonce",
              "type": "uint48"
            }
          ]
        },
        {
          "name": "spender",
          "type": "address"
        },
        {
          "name": "sigDeadline",
          "type": "uint256"
        }
      ]
    },
    {
      "name": "signature",
      "type": "bytes"
    }
  ],
  "name": "PERMIT2_PERMIT_BATCH",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
  "type": "function"
},
{
  "inputs": [
    {
      "name": "permitBatchTransferFrom",
      "type": "tuple",
      "components": [
        {
          "name": "permitted",
          "type": "tuple[]",
          "components": [
            {
              "name": "token",
              "type": "address"
            },
            {
              "name": "amount",
              "type": "uint256"
            }
          ]
        },
        {
          "name": "nonce",
          "type": "uint256"
        },
        {
          "name": "deadline",
          "type": "uint256"
        }
      ]
    },
    {
      "name": "signature",
      "type": "bytes"
    }
  ],
  "name": "PERMIT2_TRANSFER_FROM_BATCH",
  "outputs": [],
  "payable": false,
  "stateMutability": "nonpayable",
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
