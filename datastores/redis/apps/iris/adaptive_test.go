package iris_redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func (r *IrisRedisTestSuite) TestSetMetricLatencyTDigest() {
	err := IrisRedisClient.SetMetricLatencyTDigest(context.Background(), 1, "fooTable", "foo", 10.0)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestGetMetricLatencyTDigest() {
	quantileVal, sc, err := IrisRedisClient.GetMetricPercentile(context.Background(), 1, "fooTable", "foo", 0.5)
	r.NoError(err)
	r.NotEmpty(quantileVal)
	r.NotEmpty(sc)

	fmt.Println(quantileVal, sc)
}

func (r *IrisRedisTestSuite) TestGetAdaptiveEndpointByPriorityScoreAndInsertIfMissing() {
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "foo",
		},
		LatencyQuartilePercentageRank: 0,
		Latency:                       0,
		Metric:                        "",
		MetricSampleCount:             0,
	}
	err := IrisRedisClient.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(context.Background(), &tableStats)
	r.NoError(err)
}

func (r *IrisRedisTestSuite) TestSetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage() {
	tableStats := StatTable{
		OrgID:     1,
		TableName: "fooTable",
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: "foo",
		},
		LatencyQuartilePercentageRank: 0,
		Latency:                       0,
		Metric:                        "",
		MetricSampleCount:             0,
	}
	err := IrisRedisClient.SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(context.Background(), &tableStats)
	r.NoError(err)
}
