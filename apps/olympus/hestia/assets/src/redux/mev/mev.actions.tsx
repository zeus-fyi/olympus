
export interface MevState {
    bundles: FlashbotsCallBundle[];
    callBundles: [];
}

export type TraderInfoType = { [key: string]: { totalTxFees: number } };

export interface FlashbotsCallBundleResult {
    coinbaseDiff: string;      // "2717471092204423",
    ethSentToCoinbase: string; // "0",
    fromAddress: string;       // "0x37ff310ab11d1928BB70F37bC5E3cf62Df09a01c",
    gasFees: string;           // "2717471092204423",
    gasPrice: string;          // "43000001459",
    gasUsed: number;           // 63197,
    toAddress: string;         // "0xdAC17F958D2ee523a2206206994597C13D831ec7",
    txHash: string;            // "0xe2df005210bdc204a34ff03211606e5d8036740c686e9fe4e266ae91cf4d12df",
    value: string;             // "0x"
    error?: string;            // Optional because it may not always be present
    revert?: string;           // Optional because it may not always be present
}

export interface FlashbotsCallBundle {
    eventID: string,
    submissionTime: string,
    bundleHash: string,
    results: FlashbotsCallBundleResult[]
    traderInfo: TraderInfoType,
    revenue: number,
    totalCost: number,
    totalGasCost: number,
    profit: number,
}