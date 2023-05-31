import factoryArtifact from "@uniswap/v2-core/build/UniswapV2Factory.json";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {BigNumber, BigNumberish, Contract, ContractFactory, utils} from "ethers";
import pepeTokenArtifact from "../artifacts/contracts/Pepe.sol/PepeToken.json";
import wethArtifact from "@uniswap/v2-periphery/build/WETH9.json";
import routerArtifact from "@uniswap/v2-periphery/build/UniswapV2Router02.json";

const { ethers, waffle } = require("hardhat");
// @ts-ignore
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { deployUniswap, deployTokens } = require("../utils/helpers");
const hre = require("hardhat");

const wEth = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2";
const pepe = "0x6982508145454Ce325dDbE47a25d4ec3d2311933";

use(solidity);

describe("RawSwap", function () {
    let deployer: SignerWithAddress, other: SignerWithAddress;
    // @ts-ignore
    let uniswap: Contract, pairAddress, UniswapV2Factory;
    let swap: Contract, pair: Contract, weth: Contract, pepe: Contract, router: Contract, factory: Contract


    beforeEach(async function () {
        [deployer, other] = await ethers.getSigners();
        const Factory = new hre.ethers.ContractFactory(factoryArtifact.abi, factoryArtifact.bytecode, deployer);
        const factory = await Factory.deploy(deployer.address);
        await factory.deployed();


        const WETH = new ContractFactory(wethArtifact.abi, wethArtifact.bytecode, deployer);
        const weth = await WETH.deploy();
        await weth.deployed();

        const PEPE = new hre.ethers.ContractFactory(pepeTokenArtifact.abi, pepeTokenArtifact.bytecode, deployer);
        const initialSupply = utils.parseEther("42069");
        const pepe = await PEPE.deploy(initialSupply);
        await pepe.deployed();

        const Router = new hre.ethers.ContractFactory(routerArtifact.abi, routerArtifact.bytecode, deployer);
        const router = await Router.deploy(factory.address, weth.address);
        await router.deployed();

        pair = await factory.getPair(weth.address, pepe.address);
        const RawSwap = await ethers.getContractFactory("RawSwap");
        swap = await RawSwap.connect(deployer).deploy();
        await swap.deployed();

    });

    it("should perform swap", async function () {
        // Arrange
        const amountIn = ethers.utils.parseEther("1");
        const amountOutMin = ethers.utils.parseEther("0.5");
        // @ts-ignore
        await weth.connect(deployer).deposit({ value: amountIn });
        // @ts-ignore
        await weth.connect(deployer).approve(swap.address, amountIn);
        // @ts-ignore
        await pepe.connect(deployer).transfer(swap.address, amountIn);

        const data = ethers.utils.defaultAbiCoder.encode(
            ["address", "uint256", "uint256"],
            [pair.address, amountIn, amountOutMin]
        );
        const tx = await deployer.sendTransaction({
            to: swap.address,
            data,
            gasLimit: 500000,
        });

        // Assert
        await expect(tx).to.emit(weth.address, "Transfer").withArgs(swap.address, deployer.address, amountIn);

    });
})

export function expandTo18Decimals(n: number): BigNumber {
    return bigNumberify(n).mul(bigNumberify(10).pow(18))
}
export declare function bigNumberify(value: BigNumberish): BigNumber;
