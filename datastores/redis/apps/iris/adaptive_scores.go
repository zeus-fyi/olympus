package iris_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

/*
	get key from round robin
	if key doesn't exist in sorted set then add with a score of 1

	then if t-digest stat is available, multiply the score by the percentile rating times the scale factor
	which is [0.618, 1.618], eg ((1 * latency quartile) + 0.618) * previous score

	[0 * median] = ~0 change

	get rank of endpoint from round robin & if it's > 1, then decay it by 10% of the score

	TODO: send the endpoint metric stat to get the percentile value from

	Then get the stat count, and if > 10 samples then start using adaptive
*/

// eg. each previous score should decay by 5% of the score if > 1

const (
	DecayConstant                 = 0.95
	StatsTimeToLiveAfterLastUsage = 3 * time.Minute
)

type StatTable struct {
	OrgID              int     `json:"orgID"`
	TableName          string  `json:"tableName"`
	MemberRankScoreIn  redis.Z `json:"memberRankScoreIn"`
	MemberRankScoreOut redis.Z `json:"memberRankScoreOut"`

	LatencyQuartilePercentageRank float64 `json:"latencyQuartileRankPercentage"`
	Latency                       float64 `json:"latency,omitempty"`
	Metric                        string  `json:"metric,omitempty"`
	MetricLatencyMedian           float64 `json:"metricLatencyMedian,omitempty"`
	MetricLatencyTail             float64 `json:"metricLatencyTail,omitempty"`
	MetricSampleCount             int     `json:"metricSampleCount,omitempty"`

	Meter *iris_usage_meters.PayloadSizeMeter `json:""`
}

func (m *IrisCache) GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	if stats.TableName == "" {
		return fmt.Errorf("stats.TableName is empty")
	}
	endpoint, ok := stats.MemberRankScoreIn.Member.(string)
	if !ok {
		return fmt.Errorf("stats.MemberRankScore.Member is not a string")
	}

	endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(stats.OrgID, stats.TableName)
	pipe := m.Writer.TxPipeline()

	var percentileCmdMedian, percentileCmdTail *redis.Cmd
	var sampleCountCmd *redis.StringCmd

	if stats.Metric != "" {
		tableMetricKey := fmt.Sprintf("%d:%s:%s", stats.OrgID, stats.TableName, stats.Metric)
		pipe.Expire(ctx, tableMetricKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
		percentileCmdMedian = pipe.Do(ctx, "PERCENTILE.GET", tableMetricKey, 0.5)
		percentileCmdTail = pipe.Do(ctx, "PERCENTILE.GET", tableMetricKey, 0.9)

		metricTdigestSampleCountKey := fmt.Sprintf("%s:samples", tableMetricKey)
		sampleCountCmd = pipe.Get(ctx, metricTdigestSampleCountKey)
	}

	// adds new member if it doesn't exist with a starting score of 1
	pipe.ZAddNX(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	scoreInCmd := pipe.ZScore(ctx, endpointPriorityScoreKey, endpoint)
	minElemCmd := pipe.ZRangeWithScores(ctx, endpointPriorityScoreKey, 0, 0)
	pipe.Expire(ctx, endpointPriorityScoreKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}
	score, err := scoreInCmd.Result()
	if err != nil {
		log.Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}
	stats.MemberRankScoreIn.Score = score
	member, err := minElemCmd.Result()
	if err != nil {
		log.Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}
	if len(member) > 0 {
		stats.MemberRankScoreOut = redis.Z{Score: member[0].Score, Member: member[0].Member}
	} else {
		stats.MemberRankScoreOut = stats.MemberRankScoreIn
	}

	if stats.Metric != "" {
		if percentileCmdMedian != nil {
			val, rerr := percentileCmdMedian.Result()
			if rerr != nil {
				log.Warn().Err(rerr).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
				rerr = nil
			} else {
				if val != nil {
					stats.MetricLatencyMedian = val.(float64)
				}
			}
		}
		if percentileCmdTail != nil {
			val, rerr := percentileCmdTail.Result()
			if rerr != nil {
				log.Warn().Err(rerr).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
			} else {
				if val != nil {
					stats.MetricLatencyTail = val.(float64)
				}
			}
		}
		if sampleCountCmd != nil {
			count, rerr := sampleCountCmd.Result()
			if rerr != nil {
				log.Warn().Err(rerr).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
				rerr = nil
			} else {
				stats.MetricSampleCount, err = strconv.Atoi(count)
				if err != nil {
					stats.MetricSampleCount = 0
					err = nil
				}
			}
		}
	}
	return nil
}

func createAdaptiveEndpointPriorityScoreKey(orgID int, tableName string) string {
	return fmt.Sprintf("%d:%s:priority", orgID, tableName)
}

// SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage updates the endpoint priority score and rate usage
func (m *IrisCache) SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	endpointOut, ok := stats.MemberRankScoreOut.Member.(string)
	if !ok {
		return fmt.Errorf("endpointMember.MemberRankScoreOut is not a string")
	}
	if stats.TableName == "" {
		return fmt.Errorf("stats.TableName is empty")
	}
	rateLimiterKey := orgRateLimitTag(stats.OrgID)
	orgRequests := orgMonthlyUsageTag(stats.OrgID, time.Now().UTC().Month().String())
	endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(stats.OrgID, stats.TableName)

	scoreAdjustmentIncrMemberOut := ((stats.LatencyQuartilePercentageRank + 0.618) * stats.MemberRankScoreOut.Score) - stats.MemberRankScoreOut.Score
	pipe := m.Writer.TxPipeline()
	if stats.Meter != nil {
		_ = pipe.IncrByFloat(ctx, orgRequests, stats.Meter.ZeusResponseComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, stats.Meter.ZeusResponseComputeUnitsConsumed())
	}
	pipe.ZIncrBy(ctx, endpointPriorityScoreKey, scoreAdjustmentIncrMemberOut, endpointOut)
	if stats.MemberRankScoreIn.Score > 1 {
		stats.MemberRankScoreIn.Score *= DecayConstant
		pipe.ZAdd(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	}

	var tdigestResp *redis.Cmd
	if stats.Metric != "" && stats.Latency > 0 {
		tableMetricKey := fmt.Sprintf("%d:%s:%s", stats.OrgID, stats.TableName, stats.Metric)
		metricTdigestSampleCountKey := fmt.Sprintf("%s:samples", tableMetricKey)
		pipe.Incr(ctx, metricTdigestSampleCountKey)
		tdigestResp = pipe.Do(ctx, "PERCENTILE.MERGE", tableMetricKey, stats.Latency)
		pipe.Expire(ctx, tableMetricKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
	}
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage")
		return err
	}

	if tdigestResp != nil {
		err = tdigestResp.Err()
		if err != nil {
			log.Err(err).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage")
			return err
		}
	}
	return err
}

func (m *IrisCache) GetNextAdaptiveRoute(ctx context.Context, orgID int, rgName string, ri iris_models.RouteInfo, meter *iris_usage_meters.PayloadSizeMeter) (*StatTable, error) {
	ts := &StatTable{
		OrgID:     orgID,
		TableName: rgName,
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: ri.RoutePath,
		},
		MemberRankScoreOut:            redis.Z{},
		LatencyQuartilePercentageRank: 0,
		Latency:                       0,
		Metric:                        "testMetricName",
		MetricLatencyMedian:           0,
		MetricLatencyTail:             0,
		MetricSampleCount:             0,
		Meter:                         meter,
	}
	err := m.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(ctx, ts)
	if err != nil {
		return nil, err
	}
	return ts, nil

}
