package web3_client

//func (w *Web3Client) GetFilteredPendingMempoolTxs(ctx context.Context, mevTxMap MevSmartContractTxMap) (MevSmartContractTxMap, error) {
//	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
//	defer cancel()
//	if mevTxMap.MethodTxMap == nil {
//		mevTxMap.MethodTxMap = make(map[string]MevTx)
//	}
//	w.Dial()
//	defer w.Close()
//	mempool, err := w.Web3Actions.GetTxPoolContent(ctx)
//	if err != nil {
//		log.Ctx(ctx).Err(err).Msg("Web3Client| GetFilteredPendingMempoolTxs")
//		return mevTxMap, err
//	}
//	processedTxMap, err := ProcessMempoolTxs(ctx, mempool["pending"], mevTxMap)
//	if err != nil {
//		log.Ctx(ctx).Err(err).Msg("Web3Client| GetFilteredPendingMempoolTxs")
//		return mevTxMap, err
//	}
//	mevTxMap = processedTxMap
//	processedTxMap, err = ProcessMempoolTxs(ctx, mempool["mempool"], mevTxMap)
//	if err != nil {
//		log.Ctx(ctx).Err(err).Msg("Web3Client| GetRawPendingMempoolTxs")
//		return mevTxMap, err
//	}
//	mevTxMap = processedTxMap
//	return mevTxMap, nil
//}
