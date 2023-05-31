import {ethers} from "hardhat";
import {ContractFactory} from "ethers";
import rawSwapArtifact from "../artifacts/contracts/RawSwap.sol/RawSwap.json";

async function main() {
    const [deployer] = await ethers.getSigners();

    const RawSwap = new ContractFactory(rawSwapArtifact.abi, rawSwapArtifact.bytecode, deployer);
    const rawSwap = await RawSwap.deploy();
    await rawSwap.deployed();
    console.log("RawSwap deployed to:", rawSwap.address);
    console.log("Owner:", deployer.address);
}

main()
    .then(() => process.exit(0))
    .catch(error => {
            console.error(error);
            process.exit(1);
        }
    );
