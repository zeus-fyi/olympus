const { ethers, waffle } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { deployUniswap, deployTokens } = require("../utils/helpers");

use(solidity);

describe("Sandwich", function () {
  let deployer, other;
  let uniswap, tokens, sandwich;

  beforeEach(async function () {
    [deployer, other] = await ethers.getSigners();

    uniswap = await deployUniswap(deployer);
    tokens = await deployTokens(deployer);

    const Sandwich = await ethers.getContractFactory("Sandwich");
    sandwich = await Sandwich.connect(deployer).deploy(deployer.address);
    await sandwich.deployed();
  });

  xit("should perform a sandwich attack", async function () {
    // Arrange
    const amountIn = ethers.utils.parseEther("1");
    const amountOutMin = ethers.utils.parseEther("0.5");
    await tokens.weth.connect(deployer).deposit({ value: amountIn });
    await tokens.weth.connect(deployer).approve(sandwich.address, amountIn);
    await tokens.weth.connect(deployer).transfer(sandwich.address, amountIn);

    // Act
    const data = ethers.utils.hexConcat([
      "0x8980f11f", // recoverERC20(address)
      ethers.utils.defaultAbiCoder.encode(["address"], [tokens.weth.address]), // token
      ethers.utils.defaultAbiCoder.encode(["uint256"], [amountIn]), // amount
    ]);

    const tx = await deployer.sendTransaction({
      to: sandwich.address,
      data,
      gasLimit: 500000,
    });

    // Assert
    await expect(tx)
      .to.emit(tokens.weth, "Transfer")
      .withArgs(sandwich.address, deployer.address, amountIn);
  });
});
