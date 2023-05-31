import hre from "hardhat";
import {ContractFactory, utils} from "ethers";

import factoryArtifact from "@uniswap/v2-core/build/UniswapV2Factory.json";
import routerArtifact from "@uniswap/v2-periphery/build/UniswapV2Router02.json";
import wethArtifact from "@uniswap/v2-periphery/build/WETH9.json";
import pepeTokenArtifact from "../artifacts/contracts/Pepe.sol/PepeToken.json";

async function deployUniswap() {
  const [owner, otherAccount] = await hre.ethers.getSigners();

  const Factory = new hre.ethers.ContractFactory(factoryArtifact.abi, factoryArtifact.bytecode, owner);
  const factory = await Factory.deploy(owner.address);
  await factory.deployed();

  const WETH = new hre.ethers.ContractFactory(wethArtifact.abi, wethArtifact.bytecode, owner);
  const weth = await WETH.deploy();
  await weth.deployed();

  const Router = new hre.ethers.ContractFactory(routerArtifact.abi, routerArtifact.bytecode, owner);
  const router = await Router.deploy(factory.address, weth.address);
  await router.deployed();

  return { factory, router, weth };
}

async function deployTokens() {
  const [owner, otherAccount] = await hre.ethers.getSigners();

  const WETH = new ContractFactory(wethArtifact.abi, wethArtifact.bytecode, owner);
  const weth = await WETH.deploy();
  await weth.deployed();

  const PEPE = new hre.ethers.ContractFactory(pepeTokenArtifact.abi, pepeTokenArtifact.bytecode, owner);
  const initialSupply = utils.parseEther("42069");
  const pepe = await PEPE.deploy(initialSupply);
  await pepe.deployed();

  return { weth, pepe };
}

const UniswapV2FactoryAddress = "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f";
const UniswapV2FactoryABI = [
  "function getPair(address tokenA, address tokenB) external view returns (address pair)",
];
const UniswapV2RouterAbi = [
  "function getAmountsOut(uint amountIn, address[] memory path) public view returns (uint[] memory amounts)",
  "function swapExactTokensForTokens(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline) external",
  "function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline) external payable",
  "function swapExactTokensForETH(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline) external",
];
const IERC20_ABI = [
  "function balanceOf(address account) external view returns (uint256)",
  "function transfer(address recipient, uint256 amount) external returns (bool)",
  "function approve(address spender, uint256 amount) external returns (bool)",
  "function transferFrom(address sender, address recipient, uint256 amount) external returns (bool)",
  "function totalSupply() external view returns (uint256)",
  "function decimals() external view returns (uint8)",
];

const WETH_ABI = [
  // Some details about the ERC20 token
  "function balanceOf(address) view returns (uint)",
  // The name of the token
  "function name() view returns (string)",
  
  // The WETH deposit function
  "function deposit() public payable",
];




export { deployUniswap, deployTokens, UniswapV2FactoryAddress, UniswapV2FactoryABI, IERC20_ABI, WETH_ABI, UniswapV2RouterAbi };
