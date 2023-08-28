package iris_redis

import (
	"context"
	"math/rand"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

func (r *IrisRedisTestSuite) TestSetMetricLatencyTDigest() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "fooTestTable", "fooTestMetricName", 1.0)
	r.NoError(err)
}

func randomBetween(x, y int) int {
	if x > y {
		// Swap the numbers if x is greater than y
		x, y = y, x
	}
	return x + rand.Intn(y-x+1)
}

func (r *IrisRedisTestSuite) TestGetAdaptiveEndpointByPriorityScoreAndInsertIfMissing() {
	//for i := 0; i < 100; i++ {
	//	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "fooTestTable", "fooTestMetricName", float64(i))
	//	r.NoError(err)
	//}
	latency := float64(randomBetween(1, 100))
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "https://artemis.zeus.fyi",
		},
		LatencyQuartilePercentageRank: latency / 100.0,
		LatencyMilliseconds:           int64(latency),
		Metric:                        "fooTestMetricName",
		MetricSampleCount:             100,
	}
	err := IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)

	latency = float64(randomBetween(1, 100))
	uuidStr := uuid.New().String()
	tableStats = StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "https://zeus.fyi",
		},
		LatencyQuartilePercentageRank: latency / 100.0,
		LatencyMilliseconds:           int64(latency),
		Metric:                        "fooTestMetricName" + uuidStr,
		MetricSampleCount:             100,
	}
	err = IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestSetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage() {
	latency := float64(randomBetween(1, 100))
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "https://zeus.fyi",
		},
		LatencyQuartilePercentageRank: latency / 100.0,
		LatencyMilliseconds:           int64(latency),
		Metric:                        "fooTestMetricName",
		MetricSampleCount:             100,
	}
	err := IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)

	latency = float64(randomBetween(1, 100))
	tableStatsYang := StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "https://artemis.zeus.fyi",
		},
		LatencyQuartilePercentageRank: latency / 100.0,
		LatencyMilliseconds:           int64(latency),
		Metric:                        "fooTestMetricName",
		MetricSampleCount:             100,
	}
	err = IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStatsYang)
	r.NoError(err)

	tableStats2 := StatTable{
		OrgID:                         1,
		TableName:                     "fooTestTable",
		MemberRankScoreIn:             tableStats.MemberRankScoreIn,
		MemberRankScoreOut:            tableStats.MemberRankScoreOut,
		LatencyQuartilePercentageRank: rand.Float64(),
		LatencyMilliseconds:           int64(rand.Float64() * 100.0),
		Metric:                        "fooTestMetricName",
		MetricLatencyMedian:           100.0,
		MetricLatencyTail:             400.0,
		MetricSampleCount:             101,
		Meter:                         nil,
	}
	err = IrisRedisClient.SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(context.Background(), &tableStats2)
	r.NoError(err)
}
