import {ethers} from "hardhat";
import {ContractFactory} from "ethers";
import rawDawgArtifact from "../utils/Rawdawg.json";

async function main() {
    const [deployer] = await ethers.getSigners();

    const RawDawg = new ContractFactory(rawDawgArtifact.abi, rawDawgArtifact.bytecode, deployer);
    const rawDawg = await RawDawg.deploy();
    await rawDawg.deployed();
    console.log("RawDawg deployed to:", rawDawg.address);
    console.log("Owner:", deployer.address);
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    }
);
