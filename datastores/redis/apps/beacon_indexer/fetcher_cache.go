package beacon_indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	apollo_buckets "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/buckets"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
)

type FetcherCache struct {
	*redis.Client
}

func NewFetcherCache(ctx context.Context, r *redis.Client) FetcherCache {
	log.Ctx(ctx).Info().Msg("NewFetcherCache")
	log.Info().Interface("redis", r)
	return FetcherCache{r}
}

func (f *FetcherCache) SetCheckpointCache(ctx context.Context, epoch int, ttl time.Duration) (string, error) {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)

	log.Info().Msgf("SetCheckpointCache: %s", key)
	statusCmd := f.Set(ctx, fmt.Sprintf("checkpoint-epoch-%d", epoch), epoch, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("SetCheckpointCache: %s", key)
		return key, statusCmd.Err()
	}
	log.Ctx(ctx).Info().Msgf("set cache at epoch %d", epoch)
	return key, nil
}

func (f *FetcherCache) DoesCheckpointExist(ctx context.Context, epoch int) (bool, error) {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)
	log.Info().Msgf("DoesCheckpointExist: %s", key)

	chkPoint, err := f.Get(ctx, key).Int()
	switch {
	case err == redis.Nil:
		fmt.Println("DoesCheckpointExist: key does not exist")
		return chkPoint == epoch, nil
	case err != nil:
		fmt.Println("DoesCheckpointExist: Get failed", err)
		log.Err(err).Msgf("DoesCheckpointExist: %s", key)
	case chkPoint == 0:
		fmt.Println("value is empty")
	}
	return chkPoint == epoch, err
}

func (f *FetcherCache) DeleteCheckpoint(ctx context.Context, epoch int) error {
	key := fmt.Sprintf("checkpoint-epoch-%d", epoch)
	log.Info().Msgf("DeleteCheckpoint: %s", key)

	err := f.Del(ctx, key)
	if err != nil {
		log.Err(err.Err()).Msgf("DeleteCheckpoint: %s", key)
		return err.Err()
	}
	return err.Err()
}

func (f *FetcherCache) MarshalBinary(vbe beacon_api.ValidatorBalances) ([]byte, error) {
	return json.Marshal(vbe)
}

func (f *FetcherCache) SetBalanceCache(ctx context.Context, epoch int, vbe beacon_api.ValidatorBalances, ttl time.Duration) (string, error) {
	key := fmt.Sprintf("validator-balance-epoch-%d", epoch)

	log.Info().Msgf("SetBalanceCache: %s", key)
	bin, err := f.MarshalBinary(vbe)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("SetBalanceCache: %s", key)
		return key, err
	}
	err = apollo_buckets.UploadBalancesAtEpoch(ctx, key, bin)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("SetBalanceCache: UploadBalancesAtEpoch %s", key)
	}
	statusCmd := f.Set(ctx, key, bin, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("SetBalanceCache: %s", key)
		return key, statusCmd.Err()
	}
	log.Ctx(ctx).Info().Msgf("SetBalanceCache at epoch %d", epoch)
	return key, nil
}

func (f *FetcherCache) UnmarshalBinary(data []byte) (beacon_api.ValidatorBalances, error) {
	// convert data to yours, let's assume its json data
	vbe := beacon_api.ValidatorBalances{}
	err := json.Unmarshal(data, &vbe)
	return vbe, err
}

func (f *FetcherCache) GetBalanceCache(ctx context.Context, epoch int) (beacon_api.ValidatorBalances, error) {
	key := fmt.Sprintf("validator-balance-epoch-%d", epoch)
	log.Info().Msgf("SetBalanceCache: %s", key)
	emptyVbe := beacon_api.ValidatorBalances{}
	var bytes []byte
	err := f.Get(ctx, key).Scan(&bytes)
	switch {
	case err == redis.Nil:
		fmt.Println("GetBalanceCache: key does not exist in redis, checking s3 bucket")
		bytes, err = apollo_buckets.DownloadBalancesAtEpoch(ctx, key)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("GetBalanceCache: no cache found in s3 bucket")
			return emptyVbe, nil
		}
		s := string(bytes)
		bytes, err = json.Marshal(s)
	case err != nil:
		log.Err(err).Msgf("GetBalanceCache Get failed: %s", key)
	}
	vbe, err := f.UnmarshalBinary(bytes)
	if err != nil {
		log.Err(err).Msgf("GetBalanceCache Unmarshalling failed: %s", key)
		return emptyVbe, err
	}
	log.Ctx(ctx).Info().Msgf("GetBalanceCache had cache at epoch %d", epoch)
	return vbe, nil
}
