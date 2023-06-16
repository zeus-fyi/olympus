# Uniswap Test Environment

This repository contains a Hardhat project setup that deploys a basic Uniswap V2 environment on the Ethereum network. It includes deployment of WETH, Uniswap V2 Factory, Uniswap V2 Router, a token pair, and a custom ERC20 token named PepeToken. The project also provides a script to add liquidity to the Uniswap V2 pair.

## Prerequisites

Before you begin, make sure you have Node.js and npm installed on your machine. If you haven't installed Hardhat yet, you can do so by running:

```bash
npm install -g hardhat
```

You will also need a local Ethereum network for development, such as Hardhat Network or Ganache.

## Setup
Clone the repository:
```bash
git clone https://github.com/your-repo/uniswap-test-env.git
cd uniswap-test-env
npm install
```

## Deploy the contracts
1. start local node:
```
npx hardhat node
```
2. deploy
```
npx hardhat run scripts/deploy_uniswap.js  --network localhost
```
On the first deployment the addresses of the tokens and pools will be the following:
```
Uniswap V2 Factory deployed to: 0xa195ACcEB1945163160CD5703Ed43E4f78176a54
WETH deployed to: 0x6212cb549De37c25071cF506aB7E115D140D9e42
Uniswap V2 Router deployed to: 0x6F9679BdF5F180a139d01c598839a5df4860431b
Uniswap V2 Pair deployed to: 0xf4AE7E15B1012edceD8103510eeB560a9343AFd3
PEPE deployed to: 0x0bF7dE8d71820840063D4B8653Fd3F0618986faF
Pair address: 0xD8E07746cea93F313e67C4931c8A224eeAAf875b
Added liquidity to the pair
Reserves for token 0 (in PEPE): 10000.0
Reserves for token 1 (in WETH): 1000.0
Total supply (in LP tokens): 3162.277660168379331998
```


You can interact with the contracts using hardhat tasks or the console
```
npx hardhat console
```
