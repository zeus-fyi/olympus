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
	pipe.Expire(ctx, metricTdigestKey, 15*time.Second)
	pipe.Expire(ctx, metricTdigestSampleCount, 15*time.Second)

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

func (m *IrisCache) GetMetricPercentile(ctx context.Context, orgID int, tableName, metricName string, percentile float64) (float64, int64, error) {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)
	metricTdigestSampleCount := fmt.Sprintf("%s:samples", metricTdigestKey)
	// Use Redis pipeline to perform both operations
	pipe := m.Writer.Pipeline()

	percentileCmd := pipe.Do(ctx, "PERCENTILE.GET", metricTdigestKey, percentile)
	sampleCountCmd := pipe.Get(ctx, metricTdigestSampleCount)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("GetMetricPercentile")
		return 0, 0, err
	}

	resp, err := percentileCmd.Result()
	if err != nil {
		log.Err(err).Msgf("GetMetricPercentile: percentileCmd %v", resp)
		return 0, 0, err
	}
	perc, ok := resp.(float64)
	if !ok {
		perc = 0.0
	}
	sampleCount, err := sampleCountCmd.Int64()
	if err != nil {
		log.Err(err).Msgf("GetMetricPercentile: sampleCountCmd")
		return 0, 0, err
	}

	return perc, sampleCount, nil
}
