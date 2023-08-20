package iris_redis

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

func (r *IrisRedisTestSuite) TestSetMetricLatencyTDigest() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "fooTestTable", "fooTestMetricName", 10.0)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetMetricLatencyTDigest() {
	quantileVal, sc, err := IrisRedisClient.GetMetricPercentile(context.Background(), 1, "fooTestTable", "fooTestMetricName", 0.5)
	r.NoError(err)
	r.NotEmpty(quantileVal)
	r.NotEmpty(sc)

	fmt.Println(quantileVal, sc)
}

func (r *IrisRedisTestSuite) TestGetAdaptiveEndpointByPriorityScoreAndInsertIfMissing() {
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "fooTest",
		},
		LatencyQuartilePercentageRank: 0,
		Latency:                       0,
		Metric:                        "fooTestMetricName",
		MetricSampleCount:             0,
	}
	err := IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)

	uuidStr := uuid.New().String()
	tableStats = StatTable{
		OrgID:     1,
		TableName: "fooTestTable" + uuidStr,
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "fooTest" + uuidStr,
		},
		LatencyQuartilePercentageRank: 0,
		Latency:                       0,
		Metric:                        "fooTestMetricName" + uuidStr,
		MetricSampleCount:             0,
	}
	err = IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestSetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage() {
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTestTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "fooTest" + uuid.New().String(),
		},
		MemberRankScoreOut: redis.Z{
			Score:  1,
			Member: "fooTest" + uuid.New().String(),
		},
		LatencyQuartilePercentageRank: .5,
		Latency:                       float64(rand.Int()),
		Metric:                        "fooTestMetricName",
		MetricLatencyMedian:           .5,
		MetricLatencyTail:             1,
		MetricSampleCount:             10,
		Meter:                         nil,
	}
	err := IrisRedisClient.SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(context.Background(), &tableStats)
	r.NoError(err)
}
