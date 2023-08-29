package iris_redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type TableMetricsSummary struct {
	TableName string
	Routes    []redis.Z
	Metrics   map[string]TableMetric
}

type TableMetric struct {
	SampleCount       int
	MetricPercentiles []MetricSample
}

type MetricSample struct {
	Percentile float64
	Latency    float64
}

type MetricSampleWrapper struct {
	MetricSample
	RedisResult      *redis.Cmd
	RedisCountResult *redis.StringCmd
}

func (m *MetricSampleWrapper) GetMetricSample() MetricSample {
	val, err := m.RedisResult.Float64()
	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msg("GetMetricSample: failed to get percentile for metric")
		return MetricSample{}
	}
	m.Latency = val
	return m.MetricSample
}

func (m *IrisCache) GetPriorityScoresAndTdigestMetrics(ctx context.Context, orgID int, rgName string) (TableMetricsSummary, error) {
	pipe := m.Reader.TxPipeline()
	endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(orgID, rgName)
	tblMetricsSetKey := getTableMetricSetKey(orgID, rgName)

	// SMembers fetches all set members
	tblMetricsCmd := pipe.SMembers(ctx, tblMetricsSetKey)
	routesWithScoresCmd := pipe.ZRangeWithScores(ctx, endpointPriorityScoreKey, 0, -1)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("GetPriorityScoresAndTdigestMetrics: failed to execute pipeline in GetPriorityScoresAndTdigestMetrics")
		return TableMetricsSummary{}, err
	}

	tblMetrics, err := tblMetricsCmd.Result()
	if err != nil && err != redis.Nil {
		log.Err(err).Msgf("GetPriorityScoresAndTdigestMetrics: failed to get table metrics")
		return TableMetricsSummary{}, err
	}

	routesWithScores, err := routesWithScoresCmd.Result()
	if err != nil && err != redis.Nil {
		log.Err(err).Msgf("GetPriorityScoresAndTdigestMetrics: failed to get routes with scores")
		return TableMetricsSummary{}, err
	}

	pipe = m.Reader.TxPipeline()
	metricPipeline := make(map[string][]MetricSampleWrapper, len(tblMetrics))
	metricCountPipeline := make(map[string]MetricSampleWrapper, len(tblMetrics))

	for _, tbm := range tblMetrics {
		histogramBins := 7
		metricPipeline[tbm] = make([]MetricSampleWrapper, histogramBins)

		metricTdigestSampleCountKey := getMetricTdigestMetricSamplesKey(orgID, rgName, tbm)
		sampleCountCmd := pipe.Get(ctx, metricTdigestSampleCountKey)
		metricCountPipeline[tbm] = MetricSampleWrapper{
			RedisCountResult: sampleCountCmd,
		}
		metricKey := getTableMetricKey(orgID, rgName, tbm)
		for j := 0; j < histogramBins; j++ {
			percentile := 0.0
			switch j {
			case 0:
				percentile = 0.0
			case 1:
				percentile = 0.25
			case 2:
				percentile = 0.5
			case 3:
				percentile = 0.75
			case 4:
				percentile = 0.9
			case 5:
				percentile = 0.95
			case 6:
				percentile = 0.99
			}

			mw := MetricSampleWrapper{
				MetricSample: MetricSample{
					Percentile: percentile,
				},
				RedisResult: pipe.Do(ctx, "PERCENTILE.GET", metricKey, percentile),
			}
			metricPipeline[tbm][j] = mw
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Err(err).Msgf("GetPriorityScoresAndTdigestMetrics: failed to execute pipeline in GetPriorityScoresAndTdigestMetrics")
		return TableMetricsSummary{}, err
	}
	ts := TableMetricsSummary{
		TableName: rgName,
		Routes:    routesWithScores,
		Metrics:   make(map[string]TableMetric),
	}
	for metricName, tbm := range metricPipeline {
		count, cerr := metricCountPipeline[metricName].RedisCountResult.Int()
		if cerr != nil {
			log.Err(cerr).Msgf("GetPriorityScoresAndTdigestMetrics: failed to get count for metric %s", metricName)
			continue
		}
		mp := make([]MetricSample, len(tbm))
		for i, mw := range tbm {
			mv := mw.GetMetricSample()
			mp[i] = mv
		}
		ts.Metrics[metricName] = TableMetric{
			SampleCount:       count,
			MetricPercentiles: mp,
		}
	}
	return ts, nil
}
