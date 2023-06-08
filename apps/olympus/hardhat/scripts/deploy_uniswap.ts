import {ethers} from "hardhat";
import {Contract, ContractFactory, utils} from "ethers";
import sandwichArtifact from "../artifacts/contracts/Sandwich.sol/Sandwich.json";

import factoryArtifact from "@uniswap/v2-core/build/UniswapV2Factory.json";
import routerArtifact from "@uniswap/v2-periphery/build/UniswapV2Router02.json";
import wethArtifact from "@uniswap/v2-periphery/build/WETH9.json";
import pairArtifact from "@uniswap/v2-core/build/UniswapV2Pair.json";


const { parseEther } = utils;

async function main() {
    const [deployer] = await ethers.getSigners();

    const Factory = new ContractFactory(factoryArtifact.abi, factoryArtifact.bytecode, deployer);
    const factory = await Factory.deploy(deployer.address);
    await factory.deployed();
    console.log("Uniswap V2 Factory deployed to:", factory.address);

    console.log(factoryArtifact.bytecode)

    const WETH = new ContractFactory(wethArtifact.abi, wethArtifact.bytecode, deployer);
    const weth = await WETH.deploy();
    await weth.deployed();
    console.log("WETH deployed to:", weth.address);

    const Router = new ContractFactory(routerArtifact.abi, routerArtifact.bytecode, deployer);
    const router = await Router.deploy(factory.address, weth.address);
    await router.deployed();
    console.log("Uniswap V2 Router deployed to:", router.address);

    const Pair = new ContractFactory(pairArtifact.abi, pairArtifact.bytecode, deployer);
    const pair = await Pair.deploy();
    await pair.deployed();
    console.log("Uniswap V2 Pair deployed to:", pair.address);

    const PEPE = await ethers.getContractFactory("PepeToken", deployer);
    const pepeInitalSupply = parseEther("42069");
    const pepe = await PEPE.deploy(pepeInitalSupply);
    await pepe.deployed();
    console.log("PEPE deployed to:", pepe.address);

    // Mint tokens to deployer for pepe
    const pepeMintAmount = parseEther("10000");
    await pepe.mint(deployer.address, pepeMintAmount);

    // Create pair
    await factory.createPair(pepe.address, weth.address);

    // Get pair address
    const pairAddress = await factory.getPair(pepe.address, weth.address);
    console.log("Pair address:", pairAddress);

    // Add liquidity
    const pepeLiquidityAmount = parseEther("10000");
    const wethLiquidityAmount = parseEther("1000");
    await pepe.approve(router.address, pepeLiquidityAmount);
    await weth.approve(router.address, wethLiquidityAmount);

    const deadline = Math.floor(Date.now() / 1000) + 60 * 20; // 20 minutes from the current Unix time
    await router.addLiquidityETH(
        pepe.address,
        pepeLiquidityAmount,
        pepeLiquidityAmount,
        wethLiquidityAmount,
        deployer.address,
        deadline,
        { value: wethLiquidityAmount }
    );
    console.log("Added liquidity to the pair");


    const pairContract = new Contract(pairAddress, pairArtifact.abi, deployer);
    const totalSupply = await pairContract.totalSupply();
    const reserves = await pairContract.getReserves();

    console.log("Reserves for token 0 (in PEPE):", ethers.utils.formatUnits(reserves._reserve0, 18), " | ", reserves._reserve0.toString());
    console.log("Reserves for token 1 (in WETH):", ethers.utils.formatEther(reserves._reserve1), " | ", reserves._reserve1.toString());
    console.log("Total supply (in LP tokens):", ethers.utils.formatUnits(totalSupply, 18), " | ", totalSupply.toString());

    // Deploy sandwich contract
    const Sandwich = new ContractFactory(sandwichArtifact.abi, sandwichArtifact.bytecode, deployer);
    const sandwich = await Sandwich.deploy(deployer.address);
    await sandwich.deployed();
    console.log("Sandwich deployed to:", sandwich.address);

    

}


main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    }
    );
