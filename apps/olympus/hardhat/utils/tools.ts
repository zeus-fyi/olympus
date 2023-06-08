import {ethers} from "hardhat";
import * as dotenv from "dotenv";
import fetch from "node-fetch";

dotenv.config();

// @ts-ignore
async function getContractFromAddress(address: string): Promise<ethers.Contract | undefined> {
  const apiKey = process.env.ETHERSCAN_API_KEY;

  if (apiKey === undefined) {
    throw new Error("ETHERSCAN_API_KEY not set");
  }
  const url = `https://api.etherscan.io/api?module=contract&action=getabi&address=${address}&apikey=${apiKey}`;

  try {
        const response = await fetch(url);
        const data = await response.json();
        const [signer] = await ethers.getSigners();

        if (data.status === "0") {
            throw new Error(data.result);
        } else if (data.status === "1" && data.message === "OK") {
            const abi = JSON.parse(data.result);
            const provider = new ethers.providers.JsonRpcProvider();
            const contract = new ethers.Contract(address, abi, provider).connect(signer);
            return contract;
        }
    } catch (error) {
        // @ts-ignore
      throw new Error(error);
    }
  }

  export { getContractFromAddress };