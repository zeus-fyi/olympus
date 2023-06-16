// In hardhat.config.js
import {HardhatUserConfig} from "hardhat/config";
import "@nomiclabs/hardhat-waffle";
import "@nomiclabs/hardhat-ethers";
// import "hardhat-typechain"


const config: HardhatUserConfig = {
  solidity: {
    compilers: [
      { version: "0.5.16" },
      { version: "0.6.6" },
      { version: "0.6.0" },
      { version: "0.7.6" },
      { version: "0.8.0" },
      { version: "0.8.17" },
    ],
    settings: {
      allowUnlimitedContractSize: true,
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  defaultNetwork: "localhost",
  networks: {
    hardhat: {
      chainId: 1,
    },
    localhost: {
      url: "http://127.0.0.1:8545",
      gas: 8000000,
      gasPrice: 8000000000,
      blockGasLimit: 0x1fffffffffffff,
    },
  },
};

export default config;