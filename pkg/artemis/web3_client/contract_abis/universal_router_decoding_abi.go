package artemis_oly_contract_abis

const UniversalRouterAbi = `[{"inputs":[{"components":[{"internalType":"address","name":"permit2","type":"address"},{"internalType":"address","name":"weth9","type":"address"},{"internalType":"address","name":"seaport","type":"address"},{"internalType":"address","name":"nftxZap","type":"address"},{"internalType":"address","name":"x2y2","type":"address"},{"internalType":"address","name":"foundation","type":"address"},{"internalType":"address","name":"sudoswap","type":"address"},{"internalType":"address","name":"nft20Zap","type":"address"},{"internalType":"address","name":"cryptopunks","type":"address"},{"internalType":"address","name":"looksRare","type":"address"},{"internalType":"address","name":"routerRewardsDistributor","type":"address"},{"internalType":"address","name":"looksRareRewardsDistributor","type":"address"},{"internalType":"address","name":"looksRareToken","type":"address"},{"internalType":"address","name":"v2Factory","type":"address"},{"internalType":"address","name":"v3Factory","type":"address"},{"internalType":"bytes32","name":"pairInitCodeHash","type":"bytes32"},{"internalType":"bytes32","name":"poolInitCodeHash","type":"bytes32"}],"internalType":"struct RouterParameters","name":"params","type":"tuple"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"ContractLocked","type":"error"},{"inputs":[],"name":"ETHNotAccepted","type":"error"},{"inputs":[{"internalType":"uint256","name":"commandIndex","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"}],"name":"ExecutionFailed","type":"error"},{"inputs":[],"name":"FromAddressIsNotOwner","type":"error"},{"inputs":[],"name":"InsufficientETH","type":"error"},{"inputs":[],"name":"InsufficientToken","type":"error"},{"inputs":[],"name":"InvalidBips","type":"error"},{"inputs":[{"internalType":"uint256","name":"commandType","type":"uint256"}],"name":"InvalidCommandType","type":"error"},{"inputs":[],"name":"InvalidOwnerERC1155","type":"error"},{"inputs":[],"name":"InvalidOwnerERC721","type":"error"},{"inputs":[],"name":"InvalidPath","type":"error"},{"inputs":[],"name":"InvalidReserves","type":"error"},{"inputs":[],"name":"LengthMismatch","type":"error"},{"inputs":[],"name":"NoSlice","type":"error"},{"inputs":[],"name":"SliceOutOfBounds","type":"error"},{"inputs":[],"name":"SliceOverflow","type":"error"},{"inputs":[],"name":"ToAddressOutOfBounds","type":"error"},{"inputs":[],"name":"ToAddressOverflow","type":"error"},{"inputs":[],"name":"ToUint24OutOfBounds","type":"error"},{"inputs":[],"name":"ToUint24Overflow","type":"error"},{"inputs":[],"name":"TransactionDeadlinePassed","type":"error"},{"inputs":[],"name":"UnableToClaim","type":"error"},{"inputs":[],"name":"UnsafeCast","type":"error"},{"inputs":[],"name":"V2InvalidPath","type":"error"},{"inputs":[],"name":"V2TooLittleReceived","type":"error"},{"inputs":[],"name":"V2TooMuchRequested","type":"error"},{"inputs":[],"name":"V3InvalidAmountOut","type":"error"},{"inputs":[],"name":"V3InvalidCaller","type":"error"},{"inputs":[],"name":"V3InvalidSwap","type":"error"},{"inputs":[],"name":"V3TooLittleReceived","type":"error"},{"inputs":[],"name":"V3TooMuchRequested","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"RewardsSent","type":"event"},{"inputs":[{"internalType":"bytes","name":"looksRareClaim","type":"bytes"}],"name":"collectRewards","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"commands","type":"bytes"},{"internalType":"bytes[]","name":"inputs","type":"bytes[]"}],"name":"execute","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes","name":"commands","type":"bytes"},{"internalType":"bytes[]","name":"inputs","type":"bytes[]"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"execute","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"uint256[]","name":"","type":"uint256[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"},{"internalType":"bytes","name":"","type":"bytes"}],"name":"onERC1155BatchReceived","outputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"bytes","name":"","type":"bytes"}],"name":"onERC1155Received","outputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"bytes","name":"","type":"bytes"}],"name":"onERC721Received","outputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"int256","name":"amount0Delta","type":"int256"},{"internalType":"int256","name":"amount1Delta","type":"int256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"uniswapV3SwapCallback","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`

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
      "name": "amount",
      "type": "uint256"
    },
    {
      "name": "nonce",
      "type": "uint256"
    },
    {
      "name": "deadline",
      "type": "uint256"
    },
    {
      "name": "to",
      "type": "address"
    },
    {
      "name": "requestedAmount",
      "type": "uint256"
    },
    {
      "name": "owner",
      "type": "address"
    },
    {
      "name": "signature",
      "type": "bytes"
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
		"name": "batchDetails",
		"type": "tuple[]",
		"components": [
			{
				"name": "from",
				"type": "address"
			},
			{
				"name": "to",
				"type": "address"
			},
			{
				"name": "amount",
				"type": "uint160"
			},
			{
				"name": "token",
				"type": "address"
			}
		]
	}
],
"name": "PERMIT2_TRANSFER_FROM_BATCH",
"outputs": [],
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
