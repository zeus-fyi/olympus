package iris_redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type TableMetricsSummary struct {
	TableName string                 `json:"tableName"`
	Routes    []redis.Z              `json:"routes"`
	Metrics   map[string]TableMetric `json:"metrics"`
}

type TableMetric struct {
	SampleCount       int              `json:"sampleCount"`
	RedisSampleCount  *redis.StringCmd `json:"-"`
	MetricPercentiles []MetricSample   `json:"metricPercentiles"`
}

type MetricSample struct {
	Percentile  float64    `json:"percentile"`
	Latency     float64    `json:"latency"`
	RedisResult *redis.Cmd `json:"-"`
}

func (m *IrisCache) GetPriorityScoresAndTdigestMetrics(ctx context.Context, orgID int, rgName string) (TableMetricsSummary, error) {
	orgID = 1
	rgName = "fooTestTable"
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
	ts := TableMetricsSummary{
		TableName: rgName,
		Routes:    routesWithScores,
		Metrics:   make(map[string]TableMetric),
	}

	for _, tbm := range tblMetrics {
		histogramBins := 8
		metricKey := getTableMetricKey(orgID, rgName, tbm)

		metricTdigestSampleCountKey := getMetricTdigestMetricSamplesKey(orgID, rgName, tbm)
		tm := TableMetric{
			RedisSampleCount:  pipe.Get(ctx, metricTdigestSampleCountKey),
			MetricPercentiles: make([]MetricSample, histogramBins),
		}
		for j := 0; j < histogramBins; j++ {
			percentile := 0.0
			switch j {
			case 0:
				percentile = 0.1
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
			case 7:
				percentile = 1.0
			}
			tm.MetricPercentiles[j].Percentile = percentile
			tm.MetricPercentiles[j].RedisResult = pipe.Do(ctx, "PERCENTILE.GET", metricKey, percentile)
		}
		ts.Metrics[tbm] = tm
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Err(err).Msgf("GetPriorityScoresAndTdigestMetrics: failed to execute pipeline in GetPriorityScoresAndTdigestMetrics")
		return TableMetricsSummary{}, err
	}

	for metricKey, metricsRedisWrapper := range ts.Metrics {
		count, cerr := metricsRedisWrapper.RedisSampleCount.Result()
		if cerr != nil && cerr != redis.Nil {
			log.Err(cerr).Msgf("GetPriorityScoresAndTdigestMetrics: failed to get sample count for metric %s", metricKey)
			continue
		}
		tmp := ts.Metrics[metricKey]
		ci, serr := strconv.Atoi(count)
		if serr == nil {
			tmp.SampleCount = ci
		}
		for i, item := range tmp.MetricPercentiles {
			val, rerr := item.RedisResult.Result()
			if rerr != nil && rerr != redis.Nil {
				log.Err(rerr).Msgf("GetPriorityScoresAndTdigestMetrics: failed to get percentile %f for metric %s", item.Percentile, metricKey)
				continue
			}
			fv, ok := val.(float64)
			if ok {
				tmp.MetricPercentiles[i].Latency = fv
			}
		}
		ts.Metrics[metricKey] = tmp
	}
	return ts, nil
}
