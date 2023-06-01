
# assumes you have an ENV variable called MAINNET_RPC_URL
mainnet.fork:
	npx hardhat node --fork ${MAINNET_RPC_URL}

BLOCK_NUMBER ?= 17326550
mainnet.replay.fork:
	npx hardhat node --fork ${MAINNET_RPC_URL} --fork-block-number ${BLOCK_NUMBER}

start.local:
	npx hardhat node

deploy.raw.swap:
	npx hardhat run scripts/deploy_rawswap.ts  --network localhost

deploy.localhost:
	npx hardhat run scripts/deploy_uniswap.ts  --network localhost

deploy.token:
	npx hardhat run scripts/deploy_token.ts  --network localhost

export MAINNET_RPC_URL="https://proud-alpha-spree.discover.quiknode.pro/175cd672179555f519dd93e21ba1f29b6069de00/"