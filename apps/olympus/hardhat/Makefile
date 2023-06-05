
# assumes you have an ENV variable called MAINNET_RPC_URL
mainnet.fork:
	npx hardhat node --fork ${MAINNET_RPC_URL}

compile:
	npx hardhat compile

BLOCK_NUMBER ?= 17326550
mainnet.replay.fork:
	npx hardhat node --fork ${MAINNET_RPC_URL} --fork-block-number ${BLOCK_NUMBER}

start.local:
	npx hardhat node

deploy.rawdawg:
	npx hardhat run scripts/deploy_rawdawg.ts  --network localhost

deploy.localhost:
	npx hardhat run scripts/deploy_uniswap.ts  --network localhost

deploy.token:
	npx hardhat run scripts/deploy_token.ts  --network localhost

