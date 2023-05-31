import {ethers, network} from "hardhat";
import * as fs from "fs";
import {TransactionRequest} from "@ethersproject/abstract-provider";
import {EthMempoolMevTx, RpcTransaction, TradeExecutionFlow,} from '../types';
import {UniswapV2RouterAbi, WETH_ABI} from "../utils/helpers";
import {cyan, green, yellow} from 'console-log-colors';
import dotenv from "dotenv";

dotenv.config();

interface LogEntry extends EthMempoolMevTx {
  txFlowPrediction: string;
}

async function resetNetworkAndGoToBlock(blockNumber: number) {
  await network.provider.request({
    method: "hardhat_reset",
    params: [
      {
        forking: {
          jsonRpcUrl: process.env.MAINNET_RPC_URL,
          blockNumber: blockNumber - 1,
        },
      },
    ],
  });
  await network.provider.request({
    method: "evm_mine",
    params: [],
  });
}

async function main() {
  const provider = ethers.provider;
  const data: string = fs.readFileSync("./test/testTxs.json", "utf-8");
  const transactions: LogEntry[] = JSON.parse(data);
  const uniswapV2RouterAddr = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D";

  for (const tx of transactions) {
    // Reset the network to the block of the transaction.
    await resetNetworkAndGoToBlock(tx.blockNumber);

    const senderAddress = tx.from;
    console.log('-----------------------------------------')
    console.log(cyan(`-- Replaying mainnet transaction ${tx.txHash}`));
    console.log(yellow(`-- Block number: ${tx.blockNumber}`))

    // Impersonate the sender's account.
    console.log(yellow(`-- Impersonating ${senderAddress}`));
    await network.provider.request({
      method: "hardhat_impersonateAccount",
      params: [senderAddress],
    });

    const signer = await ethers.getSigner(senderAddress);

    // Parse the transaction from the log.
    const tradeExecutionFlow: TradeExecutionFlow = JSON.parse(tx.txFlowPrediction);

    const rpcTransaction: RpcTransaction = tradeExecutionFlow.tx;

    let {to, gasLimit, nonce, value, input, gasPrice} = rpcTransaction;
    gasLimit = rpcTransaction.gasPrice || "21000";
    const tradeMethod = tradeExecutionFlow.trade.tradeMethod

    if (!to) console.log('to is undefined');
    if (!gasLimit) console.log('gasLimit is undefined');
    if (!nonce) console.log('nonce is undefined');
    if (!value) console.log('value is undefined');
    if (!input) console.log('input is undefined');
    if (!gasPrice) { 
        console.log('gasPrice is undefined');
        // gasPrice = "100000000000";
        gasPrice = "54000000000";
    }

    console.log(yellow(`-- gasPrice: ${ethers.BigNumber.from(gasPrice)}`))

    // log the balance weth balance of the sender
    const wethContract = await ethers.getContractAt(
        WETH_ABI,
        "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
        signer
    );

    const uniswapV2RouterContract = await ethers.getContractAt(
        UniswapV2RouterAbi,
        uniswapV2RouterAddr,
        signer
    );

    const wethBalance = await wethContract.balanceOf(senderAddress);
    const ethBalance = await provider.getBalance(senderAddress);
    console.log(yellow(`-- Method: ${tradeMethod}`))
    console.log(yellow(`-- WETH balance: ${wethBalance.toString()}`));
    console.log(yellow(`-- ETH balance: ${ethBalance.toString()}`));
    console.log(yellow(`-- Value: ${ethers.BigNumber.from(value)}`));

    
    const transaction: TransactionRequest = {
      to,
    //   gasLimit: ethers.BigNumber.from(8000000),
      value: ethers.BigNumber.from(value),
      data: input,
      gasPrice: ethers.BigNumber.from(gasPrice),
    };

    transaction.nonce = await provider.getTransactionCount(senderAddress, 'pending');


    try {
        // const gasLimitEst = await uniswapV2RouterContract.estimateGas.swapExactTokensForTokens();
        // transaction.gasLimit = gasLimitEst.add(ethers.BigNumber.from(100000));


    // Send the transaction.
        const txResponse = await signer.sendTransaction(transaction);
        const txReceipt = await provider.waitForTransaction(txResponse.hash);
        
        console.log(green(`Transaction ${txResponse.hash} mined in block ${txReceipt.blockNumber}`));
    } catch (error) {
        console.log(error);
    } finally {
        // Stop impersonating the sender.
        console.log(yellow(`-- Stopping impersonation of ${senderAddress}`));
        await network.provider.request({
            method: "hardhat_stopImpersonatingAccount",
            params: [senderAddress],
        });
        console.log('-----------------------------------------')
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
  });
