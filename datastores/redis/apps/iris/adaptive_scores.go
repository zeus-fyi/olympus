package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
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
	MetricSampleCount             int     `json:"metricSampleCount,omitempty"`

	Meter *iris_usage_meters.PayloadSizeMeter `json:""`
}

// TODO: insert sample into tigest

func (m *IrisCache) GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	endpoint, ok := stats.MemberRankScoreIn.Member.(string)
	if !ok {
		return fmt.Errorf("stats.MemberRankScore.Member is not a string")
	}

	endpointPriorityScoreKey := fmt.Sprintf("%d:%s:priority", stats.OrgID, stats.TableName)
	pipe := m.Writer.TxPipeline()
	pipe.ZAddNX(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	scoreCmd := pipe.ZScore(ctx, endpointPriorityScoreKey, endpoint)
	minElemCmd := pipe.ZRangeWithScores(ctx, endpointPriorityScoreKey, 0, 0)
	pipe.Expire(ctx, endpointPriorityScoreKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}
	score, err := scoreCmd.Result()
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
	return nil
}

// TODO, add metric table stat count + tidgest entry

func (m *IrisCache) SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	endpointOut, ok := stats.MemberRankScoreOut.Member.(string)
	if !ok {
		return fmt.Errorf("endpointMember.MemberRankScoreOut is not a string")
	}
	rateLimiterKey := orgRateLimitTag(stats.OrgID)
	orgRequests := orgMonthlyUsageTag(stats.OrgID, time.Now().UTC().Month().String())
	endpointPriorityScoreKey := fmt.Sprintf("%d:%s:priority", stats.OrgID, stats.TableName)

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
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage")
		return err
	}
	return err
}
