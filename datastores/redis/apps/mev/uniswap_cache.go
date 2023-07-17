package redis_mev

//
//func (m *MevCache) AddOrUpdatePairPricesCache(ctx context.Context, tag string, pd artemis_uniswap_pricing.UniswapPricingData, ttl time.Duration) error {
//	bin, err := m.MarshalBinary(pd)
//	if err != nil {
//		return err
//	}
//	statusCmd := m.Set(ctx, tag, bin, ttl)
//	if statusCmd.Err() != nil {
//		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("AddOrUpdateLatestBlockCache: %s", tag)
//		return statusCmd.Err()
//	}
//	return nil
//}
//
//func (m *MevCache) MarshalBinary(up artemis_uniswap_pricing.UniswapPricingData) ([]byte, error) {
//	return json.Marshal(up)
//}
//
//func (m *MevCache) UnmarshalBinary(data []byte) (artemis_uniswap_pricing.UniswapPricingData, error) {
//	pd := artemis_uniswap_pricing.UniswapPricingData{}
//	err := json.Unmarshal(data, &pd)
//	return pd, err
//}
//
//func (m *MevCache) GetPairPricesFromCacheIfExists(ctx context.Context, tag string) (artemis_uniswap_pricing.UniswapPricingData, error) {
//	pd := artemis_uniswap_pricing.UniswapPricingData{}
//	var bytes []byte
//	err := m.Get(ctx, tag).Scan(&bytes)
//	switch {
//	case err == redis.Nil:
//		return pd, fmt.Errorf("GetPairPricesFromCacheIfExists: %s", tag)
//	case err != nil:
//		log.Err(err).Msgf("GetPairPricesFromCacheIfExists Get failed: %s", tag)
//	}
//	cachedPd, err := m.UnmarshalBinary(bytes)
//	if err != nil {
//		return pd, err
//	}
//	return cachedPd, nil
//}
