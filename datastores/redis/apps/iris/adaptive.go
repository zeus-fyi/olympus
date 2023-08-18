package iris_redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
)

func (m *IrisCache) SetMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string, latencies []float64, sorted bool) error {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)

	// Convert latencies to interface slice for Redis command
	latenciesInterface := make([]interface{}, len(latencies))
	for i, v := range latencies {
		latenciesInterface[i] = v
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	// Determine whether to use PERCENTILE.MERGE or PERCENTILE.MERGESORTED based on sorted flag
	if sorted {
		pipe.Do(ctx, "PERCENTILE.MERGESORTED", metricTdigestKey, latenciesInterface)
	} else {
		pipe.Do(ctx, "PERCENTILE.MERGE", metricTdigestKey, latenciesInterface)
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetMetricLatencyTDigest")
		return err
	}
	return nil
}

func (m *IrisCache) GetMetricPercentile(ctx context.Context, orgID int, tableName, metricName string, percentile float64) (float64, error) {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)

	resp, err := m.Reader.Do(ctx, "PERCENTILE.GET", metricTdigestKey, percentile).Result()
	if err != nil {
		log.Err(err).Msgf("GetMetricPercentile")
		return 0, err
	}

	value, err := redis.Float64(resp, err)
	if err != nil {
		log.Err(err).Msgf("GetMetricPercentile conversion error")
		return 0, err
	}

	return value, nil
}
