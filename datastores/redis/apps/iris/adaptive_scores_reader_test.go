package iris_redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

// {7138983863666903883}.ethereum-mainnet:eth_blockNumber:samples
// {7138983863666903883}.ethereum-mainnet:eth_getBlockByNumber:samples
// {7138983863666903883}.ethereum-mainnet:metrics

func (r *IrisRedisTestSuite) TestReadMetrics() {
	tm, err := IrisRedisClient.GetPriorityScoresAndTdigestMetrics(context.Background(), 7138983863666903883, "ethereum-mainnet")
	r.NoError(err)
	fmt.Println(tm)
	r.Equal("ethereum-mainnet", tm.TableName)
	r.NotEmpty(tm.Metrics)
	r.NotEmpty(tm.Routes)
}

func GetMetricOffset(offset int) string {
	return fmt.Sprintf("fooTestMetricName%d", offset)
}
func (r *IrisRedisTestSuite) TestGetPriorityScoresAndTdigestMetrics() {
	pipe := IrisRedisClient.Writer.TxPipeline()

	tableName := "fooTestTable"
	routes := []iris_models.RouteInfo{
		{
			RoutePath: "https://zeus.fyi",
			Referrers: []string{"https://google.com", "https://yahoo.com"},
		},
		{
			RoutePath: "https://artemis.zeus.fyi",
			Referrers: nil,
		},
	}
	err := IrisRedisClient.AddOrUpdateOrgRoutingGroup(context.Background(), 1, tableName, routes)
	r.NoError(err)

	m1 := GetMetricOffset(1)
	m1Key := getMetricTdigestKey(1, tableName, m1)
	m2 := GetMetricOffset(2)
	m2Key := getMetricTdigestKey(1, tableName, m2)
	m3 := GetMetricOffset(3)
	m3Key := getMetricTdigestKey(1, tableName, m3)

	fmt.Println(m1Key, m2Key, m3Key)

	tblMetricsSetKey := getTableMetricSetKey(1, tableName)
	pipe.SAdd(context.Background(), tblMetricsSetKey, m1)
	pipe.SAdd(context.Background(), tblMetricsSetKey, m2)
	pipe.SAdd(context.Background(), tblMetricsSetKey, m3)
	pipe.Expire(context.Background(), tblMetricsSetKey, StatsTimeToLiveAfterLastUsage)

	pipe.IncrBy(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m1), 10)
	pipe.Expire(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m1), StatsTimeToLiveAfterLastUsage)

	pipe.IncrBy(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m2), 20)
	pipe.Expire(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m2), StatsTimeToLiveAfterLastUsage)

	pipe.IncrBy(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m3), 50)
	pipe.Expire(context.Background(), getMetricTdigestMetricSamplesKey(1, tableName, m3), StatsTimeToLiveAfterLastUsage)

	endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(1, "fooTestTable")

	ctx := context.Background()
	pipe.ZAdd(ctx, endpointPriorityScoreKey, redis.Z{
		Score:  0.8,
		Member: "https://zeus.fyi",
	})

	pipe.ZAdd(ctx, endpointPriorityScoreKey, redis.Z{
		Score:  1.2,
		Member: "https://artemis.zeus.fyi",
	})

	for i := 0; i < 10; i++ {
		v1 := float64(i)
		v2 := float64(i * 10)
		v3 := float64(i * 100)

		fmt.Println(m1Key)
		pipe.Do(ctx, "PERCENTILE.MERGE", m1Key, v1)
		pipe.Do(ctx, "PERCENTILE.MERGE", m2Key, v2)
		pipe.Do(ctx, "PERCENTILE.MERGE", m3Key, v3)
	}
	pipe.Expire(ctx, m1Key, StatsTimeToLiveAfterLastUsage)
	pipe.Expire(ctx, m2Key, StatsTimeToLiveAfterLastUsage)
	pipe.Expire(ctx, m3Key, StatsTimeToLiveAfterLastUsage)

	scoreInCmd := pipe.ZScore(ctx, endpointPriorityScoreKey, "https://zeus.fyi")
	minElemCmd := pipe.ZRangeWithScores(ctx, endpointPriorityScoreKey, 0, 0)
	pipe.Expire(ctx, endpointPriorityScoreKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
	_, err = pipe.Exec(ctx)
	//r.NoError(err)
	scoreInCmd.Result()
	minElemCmd.Result()

	tm, err := IrisRedisClient.GetPriorityScoresAndTdigestMetrics(context.Background(), 1, "fooTestTable")
	r.NoError(err)
	r.Equal("fooTestTable", tm.TableName)
	r.NotEmpty(tm.Metrics)
	r.NotEmpty(tm.Routes)

	for tbm, metric := range tm.Metrics {
		fmt.Println(tbm)
		fmt.Println(metric.SampleCount)
		for _, mp := range metric.MetricPercentiles {
			fmt.Println(mp.Latency)
			fmt.Println(mp.Percentile)
		}
	}
}
