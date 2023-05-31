import {ethers} from "hardhat";
import {ContractFactory} from "ethers";
import tokenArtifact from "../artifacts/contracts/Token.sol/Token.json";
import {parseEther} from "ethers/lib/utils";

async function main() {
    const [deployer] = await ethers.getSigners();

    const mintAmt = parseEther("10000");

    const Token = new ContractFactory(tokenArtifact.abi, tokenArtifact.bytecode, deployer);
    console.log(tokenArtifact.bytecode)
    console.log(tokenArtifact.bytecode.length)
    const tkn = await Token.deploy(mintAmt);
    await tkn.deployed();

    console.log("RawDawg deployed to:", tkn.address);
    console.log("Owner:", deployer.address);
}

main()
    .then(() => process.exit(0))
    .catch(error => {
            console.error(error);
            process.exit(1);
        }
    );
