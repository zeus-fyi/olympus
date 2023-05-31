import {BigNumber} from 'bignumber.js';

interface EthMempoolMevTx {
  txID: number;
  to: string;
  protocolNetworkID: number;
  txFlowPrediction: string;
  txHash: string;
  nonce: number;
  from: string;
  blockNumber: number;
  tx: string;
}

interface TradeExecutionFlow {
  currentBlockNumber: BigNumber;
  tx: RpcTransaction;
  // tradeMethod: string;
  // tradeParams: any;
  trade: Trade;
  initialPair: UniswapV2Pair;
  frontRunTrade: TradeOutcome;
  userTrade: TradeOutcome;
  sandwichTrade: TradeOutcome;
  sandwichPrediction: SandwichTradePrediction;
}

interface Trade {
  tradeMethod: string;
  method: any;
}

interface RpcTransaction {
  nonce: string;
  gasPrice: string;
  gasLimit: string;
  gasFeeCap?: string;
  gasTipCap?: string;
  to: string;
  value: string;
  input: string;
  from: string;
  v: string;
  r: string;
  s: string;
  hash: string;
  blockNumber?: string;
  blockHash?: string;
  transactionIndex?: string;
}

interface UniswapV2Pair {
  pairContractAddr: string;
  price0CumulativeLast: BigNumber;
  price1CumulativeLast: BigNumber;
  kLast: BigNumber;
  token0: string;
  token1: string;
  reserve0: BigNumber;
  reserve1: BigNumber;
  blockTimestampLast: BigNumber;
}

interface TradeOutcome {
  amountIn: BigNumber;
  amountInAddr: string;
  amountFees: BigNumber;
  amountOut: BigNumber;
  amountOutAddr: string;
  startReservesToken0: BigNumber;
  startReservesToken1: BigNumber;
  endReservesToken0: BigNumber;
  endReservesToken1: BigNumber;
}

interface SandwichTradePrediction {
  sellAmount: BigNumber;
  expectedProfit: BigNumber;
}

export {
  EthMempoolMevTx,
  TradeExecutionFlow,
  RpcTransaction,
  UniswapV2Pair,
  TradeOutcome,
  SandwichTradePrediction,
};
