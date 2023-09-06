package iris_redis

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (m *IrisCache) GetDetailedTableStats(ctx context.Context, orgID int, tableName, metricName string, percentile float64) (float64, int64, error) {
	metricTdigestKey := getMetricTdigestKey(orgID, tableName, metricName)
	metricTdigestSampleCount := getMetricTdigestMetricSamplesKey(orgID, tableName, metricName)
	// Use Redis pipeline to perform both operations
	pipe := m.Reader.Pipeline()

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
