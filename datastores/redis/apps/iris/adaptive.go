package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func (m *IrisCache) SetMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string, latency float64) error {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)
	metricTdigestSampleCount := fmt.Sprintf("%s:samples", metricTdigestKey)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	pipe.Incr(ctx, metricTdigestSampleCount)
	pipe.Do(ctx, "PERCENTILE.MERGE", metricTdigestKey, latency)
	pipe.Expire(ctx, metricTdigestKey, 15*time.Minute)
	pipe.Expire(ctx, metricTdigestSampleCount, 15*time.Minute)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetMetricLatencyTDigest")
		return err
	}
	return nil
}

func (m *IrisCache) DelMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string) error {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)
	metricTdigestSampleCount := fmt.Sprintf("%s:samples", metricTdigestKey)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	pipe.Del(ctx, metricTdigestSampleCount)
	pipe.Del(ctx, metricTdigestKey)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DelMetricLatencyTDigest")
		return err
	}
	return nil
}
