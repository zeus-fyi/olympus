import hre, {ethers} from "hardhat";
import {expect} from "chai";
import {BigNumberish, Contract, Signer} from "ethers";
import {IERC20_ABI} from "../utils/helpers";
import factoryArtifact from "@uniswap/v2-core/build/UniswapV2Factory.json";
import pairArtifact from "@uniswap/v2-core/build/UniswapV2Pair.json";
import routerArtifact from "@uniswap/v2-periphery/build/UniswapV2Router02.json";
import dotenv from "dotenv";
import {FactoryOptions} from "@nomiclabs/hardhat-ethers/types";

dotenv.config();
    
const wEth = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2";
const pepe = "0x6982508145454Ce325dDbE47a25d4ec3d2311933";

async function sendContractTx(
    signer: ethers.Signer,
    pairAddress: string,
    amountIn: ethers.BigNumber,
    amountOut: ethers.BigNumber,
    whichToken: number,
    contractAddress: string
) {
    const callData1 = ethers.utils.defaultAbiCoder.encode(
        ["address", "address", "uint128", "uint128", "uint8"],
        [wEth, pairAddress, amountIn, amountOut, 0]
    );

    const tx = await signer.sendTransaction({
        to: contractAddress,
        data: callData1,
    });

    const receipt = await tx.wait();
    return receipt;
}

async function getPairReserves(pairAddress: string, signer: ethers.Signer) {
    const pairContract = new ethers.Contract(pairAddress, pairArtifact.abi, signer);
    const reserves = await pairContract.getReserves();
    const totalSupply = await pairContract.totalSupply();
    console.log("Reserves for token 0 (in PEPE):", ethers.utils.formatUnits(reserves._reserve0, 18), " | ", reserves._reserve0.toString());
    console.log("Reserves for token 1 (in WETH):", ethers.utils.formatEther(reserves._reserve1), " | ", reserves._reserve1.toString());
    console.log("Total supply (in LP tokens):", ethers.utils.formatUnits(totalSupply, 18), " | ", totalSupply.toString());

    console.log("signer balance: ", ethers.utils.formatEther(await signer.getBalance()));
    return reserves;
}
   
async function instantiateContracts(pairAddress: string, signer: ethers.Signer) {
    const pairContract = await hre.ethers.getContractAt(pairArtifact.abi, pairAddress, signer);
    const pepContract = await hre.ethers.getContractAt(IERC20_ABI, pepe, signer);
    const wethContract = await hre.ethers.getContractAt(IERC20_ABI, wEth, signer);
    const routerContract = await hre.ethers.getContractAt(
      routerArtifact.abi,
      "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
      signer
    );
  
    return {
      pairContract,
      pepContract,
      wethContract,
      routerContract
    };
  }


describe("RawDawg", function () {
    let rawDawg: Contract;
    let pairAddress: string;
    let deployer, user1: Signer | FactoryOptions | undefined, user2: { address: any; getBalance: () => BigNumberish | PromiseLike<BigNumberish>; };
    let UniswapV2Factory: Contract;


    beforeEach(async function () {
        [deployer, user1, user2] = await ethers.getSigners();
        const RawDawgFactory = await ethers.getContractFactory("Rawdawg", user1);
        rawDawg = await RawDawgFactory.deploy();
        await rawDawg.deployed();

        UniswapV2Factory = await hre.ethers.getContractAt(
            factoryArtifact.abi,
            "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f"
        );
        pairAddress = await UniswapV2Factory.getPair(wEth, pepe);
        if (!pairAddress) {
            console.log("Pair doesn't exist, creating it now");
            await UniswapV2Factory.createPair(wEth, pepe);
            pairAddress = await UniswapV2Factory.getPair(wEth, pepe);
        }
        // Reset the forked node to the initial state
        await hre.network.provider.request({
            method: "hardhat_reset",
            params: [
                {
                    forking: {
                        jsonRpcUrl: process.env.MAINNET_RPC_URL,
                    },
                },
            ],
        });

    });

    it("should have pairs", async function () {
        const allPairsLength = await UniswapV2Factory.allPairsLength();
        console.log("allPairsLength: ", allPairsLength.toString());
        expect(allPairsLength).to.not.equal(0);
    });
    

    it("send first slice tx", async function () {

        const amountIn = ethers.utils.parseUnits("123", "18");
        const amountOut = ethers.utils.parseUnits("140373295133.65332", "18");

        // instantiate the pair contract
        const { pairContract, pepContract, wethContract } = await instantiateContracts(pairAddress, user1);

    
        await getPairReserves(pairAddress, user1);
        console.log("user1 PEPE balance: ", ethers.utils.formatUnits(await pepContract.balanceOf(user1.address), 18));
        console.log("user1 Eth balance: ", ethers.utils.formatEther(await  user1.getBalance()));
        console.log("sending tx...");
        const receipt1 = await sendContractTx(
            user1,
            pairAddress,
            amountIn,
            amountOut,
            1,
            rawDawg.address
        );
        expect(receipt1.status).to.equal(1);
        console.log("receipt1: ", receipt1);
        await getPairReserves(pairAddress, user1);
        console.log("user1 PEPE balance: ", ethers.utils.formatUnits(await pepContract.balanceOf(user1.address), 18));
        console.log("user1 WETH balance: ", ethers.utils.formatEther(await user1.getBalance()));

        console.log("-------------------------- 1st slice tx complete --------------------------")
    });
    
    it("Send normal swap tx", async function () {
        // instantiate the pair contract
        const { pairContract, pepContract, wethContract, routerContract } = await instantiateContracts(pairAddress, user1);


        const reserves = await getPairReserves(pairAddress, user2);
        console.log("user2 PEPE balance: ", ethers.utils.formatUnits(await pepContract.balanceOf(user2.address), 18));
        console.log("user2 ETH balance: ", ethers.utils.formatEther(await user2.getBalance())); 



        // TODO change these to parameters
        const amountIn = ethers.utils.parseEther("123");
        const amountOut = ethers.utils.parseUnits("140373295133.65332", "18");
        const deadline = Math.floor(Date.now() / 1000) + 60 * 20; // 20 minutes from the current Unix time

        // Swap Exact ETH For Tokens
        const tx = await routerContract.swapExactETHForTokens(
            amountOut,
            [wEth, pepe],
            user2.address,
            deadline,
            { value: amountIn }
        );

        const receipt = await tx.wait();
        expect(receipt.status).to.equal(1);
        // console.log("receipt: ", receipt);
        await getPairReserves(pairAddress, user2);
        console.log("user2 PEPE balance: ", ethers.utils.formatUnits(await pepContract.balanceOf(user2.address), 18));
        console.log("user2 ETH balance: ", ethers.utils.formatEther(await user2.getBalance()));
        
    });

    // it("Send normal swap and sandwich with contract", async function () {
    //     // instantiate the pair contract
    //     const pairContract = await hre.ethers.getContractAt(
    //         pairArtifact.abi,
    //         pairAddress,
    //         user2
    //     );
    //     const pepeContract = await hre.ethers.getContractAt(
    //         IERC20_ABI,
    //         pepe,
    //         user2
    //     );
    //     const wethContract = await hre.ethers.getContractAt(
    //         WETH_ABI,
    //         wEth,
    //         user2
    //     );
    //     const routerContract = await hre.ethers.getContractAt(
    //         routerArtifact.abi,
    //         "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
    //         user2
    //       });

    //     const reserves = await getPairReserves(pairAddress, user2);
    //     console.log("user2 PEPE balance: ", ethers.utils.formatUnits(await pepeContract.balanceOf(user2.address), 18));
    //     console.log("user2 ETH balance: ", ethers.utils.formatEther(await user2.getBalance()));
    
});