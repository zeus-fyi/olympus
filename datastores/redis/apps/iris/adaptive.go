package iris_redis

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (m *IrisCache) SetMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string, latency float64) error {
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	pipe.Incr(ctx, getMetricTdigestMetricSamplesKey(orgID, tableName, metricName))
	pipe.Do(ctx, "PERCENTILE.MERGE", getMetricTdigestKey(orgID, tableName, metricName), latency)
	pipe.Expire(ctx, getMetricTdigestKey(orgID, tableName, metricName), StatsTimeToLiveAfterLastUsage)
	pipe.Expire(ctx, getMetricTdigestMetricSamplesKey(orgID, tableName, metricName), StatsTimeToLiveAfterLastUsage)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetMetricLatencyTDigest")
		return err
	}
	return nil
}

func (m *IrisCache) DelMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string) error {
	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	pipe.Del(ctx, getMetricTdigestMetricSamplesKey(orgID, tableName, metricName))
	pipe.Del(ctx, getMetricTdigestKey(orgID, tableName, metricName))

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DelMetricLatencyTDigest")
		return err
	}
	return nil
}
