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
	DecayConstant                   = 0.95
	LatencyScaleFactorDefault       = 0.52
	TailPercentage                  = 0.95
	ErrorDefaultScaleFactor         = 3.0
	MinSamplesBeforeAdaptiveScoring = 20
	StatsTimeToLiveAfterLastUsage   = 60 * time.Minute
)

type StatTable struct {
	OrgID              int     `json:"orgID"`
	TableName          string  `json:"tableName"`
	MemberRankScoreIn  redis.Z `json:"memberRankScoreIn"`
	MemberRankScoreOut redis.Z `json:"memberRankScoreOut"`

	LatencyQuartilePercentageRank float64 `json:"latencyQuartileRankPercentage"`
	LatencyMilliseconds           int64   `json:"latency,omitempty"`
	Metric                        string  `json:"metric,omitempty"`
	MetricLatencyMedian           float64 `json:"metricLatencyMedian,omitempty"`
	MetricLatencyTail             float64 `json:"metricLatencyTail,omitempty"`
	MetricSampleCount             int     `json:"metricSampleCount,omitempty"`
	ScaleFactor                   float64 `json:"scaleFactor,omitempty"`

	LatencyScaleFactor float64 `json:"latencyScaleFactor,omitempty"`
	ErrorScaleFactor   float64 `json:"errorScaleFactor,omitempty"`
	DecayScaleFactor   float64 `json:"decayScaleFactor,omitempty"`

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
		tm := getTableMetricKey(stats.OrgID, stats.TableName, stats.Metric)
		pipe.Expire(ctx, tm, StatsTimeToLiveAfterLastUsage)
		percentileCmdMedian = pipe.Do(ctx, "PERCENTILE.GET", tm, 0.5)
		percentileCmdTail = pipe.Do(ctx, "PERCENTILE.GET", tm, TailPercentage)

		metricTdigestSampleCountKey := getMetricTdigestMetricSamplesKey(stats.OrgID, stats.TableName, stats.Metric)
		pipe.Expire(ctx, metricTdigestSampleCountKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
		sampleCountCmd = pipe.Get(ctx, metricTdigestSampleCountKey)

		tblMetricSet := getTableMetricSetKey(stats.OrgID, stats.TableName)
		pipe.SAdd(ctx, tblMetricSet, stats.Metric)
		pipe.Expire(ctx, tblMetricSet, StatsTimeToLiveAfterLastUsage)
	}

	// adds new member if it doesn't exist with a starting score of 1
	pipe.ZAddNX(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	scoreInCmd := pipe.ZScore(ctx, endpointPriorityScoreKey, endpoint)
	minElemCmd := pipe.ZRangeWithScores(ctx, endpointPriorityScoreKey, 0, 0)
	pipe.Expire(ctx, endpointPriorityScoreKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes
	// Execute the transaction
	latSfKey := createAdaptiveEndpointPriorityScoreLatencyScaleFactorKey(stats.OrgID, stats.TableName)
	errSfKey := createAdaptiveEndpointPriorityScoreErrorScaleFactorKey(stats.OrgID, stats.TableName)
	decaySfKey := createAdaptiveEndpointPriorityScoreDecayScaleFactorKey(stats.OrgID, stats.TableName)
	latSfCmd := pipe.Get(ctx, latSfKey)
	errSfCmd := pipe.Get(ctx, errSfKey)
	decaySfCmd := pipe.Get(ctx, decaySfKey)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}

	latSfValue, err := latSfCmd.Float64()
	if err == redis.Nil {
		latSfValue = LatencyScaleFactorDefault
		stats.LatencyScaleFactor = latSfValue
	} else if err != nil {
		log.Warn().Err(err).Msgf("Failed to get latSfKey")
	} else {
		stats.LatencyScaleFactor = latSfValue
	}

	errSfValue, err := errSfCmd.Float64()
	if err == redis.Nil {
		errSfValue = ErrorDefaultScaleFactor
	} else if err != nil {
		log.Warn().Err(err).Msgf("Failed to get errSfKey")
	} else {
		stats.ErrorScaleFactor = errSfValue
	}

	decaySfValue, err := decaySfCmd.Float64()
	if err == redis.Nil {
		decaySfValue = DecayConstant
	} else if err != nil {
		log.Warn().Err(err).Msgf("Failed to get decaySfKey")
	} else {
		stats.DecayScaleFactor = decaySfValue
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
		path, eok := member[0].Member.(string)
		if !eok {
			stats.MemberRankScoreOut = stats.MemberRankScoreIn
		}
		if len(path) != 0 {
			stats.MemberRankScoreOut = redis.Z{Score: member[0].Score, Member: member[0].Member}
		} else {
			stats.MemberRankScoreOut = stats.MemberRankScoreIn
		}
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

// SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage updates the endpoint priority score and rate usage
func (m *IrisCache) SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	if stats.TableName == "" {
		return fmt.Errorf("stats.TableName is empty")
	}

	orgRequests := getOrgMonthlyUsageKey(stats.OrgID, time.Now().UTC().Month().String())
	endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(stats.OrgID, stats.TableName)

	if stats.LatencyQuartilePercentageRank <= 0.0 {
		if stats.LatencyMilliseconds > int64(stats.MetricLatencyTail) {
			stats.LatencyQuartilePercentageRank = 1.0
		} else if stats.LatencyMilliseconds < int64(stats.MetricLatencyMedian) {
			stats.LatencyQuartilePercentageRank = 0.5
		} else {
			stats.LatencyQuartilePercentageRank = 0.75
		}
	}

	//log.Info().Float64(" stats.MetricLatencyMedian", stats.MetricLatencyMedian).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")
	//log.Info().Float64(" stats.MetricLatencyTail", stats.MetricLatencyTail).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")
	//log.Info().Int64(" stats.LatencyMilliseconds", stats.LatencyMilliseconds).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")

	/* this is being set upstream
	if resp.StatusCode >= 400 {
		tableStats.LatencyScaleFactor = tableStats.ErrorScaleFactor
	}
	*/
	rate := stats.LatencyQuartilePercentageRank + stats.LatencyScaleFactor
	// essentially this just multiplies the score by the priority rate growth
	scoreAdjustmentMemberOut := rate * stats.MemberRankScoreOut.Score
	stats.MemberRankScoreOut.Score = scoreAdjustmentMemberOut
	pipe := m.Writer.TxPipeline()

	rateLimiterKey := getOrgRateLimitKey(stats.OrgID)
	if stats.Meter != nil {
		_ = pipe.IncrByFloat(ctx, orgRequests, stats.Meter.ZeusResponseComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, stats.Meter.ZeusResponseComputeUnitsConsumed())
	}

	if stats.MetricSampleCount >= MinSamplesBeforeAdaptiveScoring {
		pipe.ZAdd(ctx, endpointPriorityScoreKey, stats.MemberRankScoreOut)
	}

	if stats.MemberRankScoreIn.Score > 1 {
		if stats.MemberRankScoreIn.Score > 100 && stats.DecayScaleFactor > 0.92 {
			stats.DecayScaleFactor = 0.92
		}
		if stats.MemberRankScoreIn.Score > 1000 && stats.DecayScaleFactor >= 0.92 {
			stats.DecayScaleFactor = 0.90
		}
		if stats.MemberRankScoreIn.Score > 10000 && stats.DecayScaleFactor >= 0.90 {
			stats.DecayScaleFactor = 0.80
		}
		if stats.MemberRankScoreIn.Score > 100000 && stats.DecayScaleFactor >= 0.80 {
			stats.DecayScaleFactor = 0.60
		}
		if stats.MemberRankScoreIn.Score > 1000000 && stats.DecayScaleFactor >= 0.60 {
			stats.DecayScaleFactor = 0.50
		}
		if stats.MemberRankScoreIn.Score > 10000000 && stats.DecayScaleFactor >= 0.50 {
			stats.DecayScaleFactor = 0.40
		}
		stats.MemberRankScoreIn.Score *= stats.DecayScaleFactor
		pipe.ZAdd(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	}

	var tdigestResp *redis.Cmd
	if stats.Metric != "" && stats.LatencyMilliseconds > 0 {
		tableMetricKey := getMetricTdigestKey(stats.OrgID, stats.TableName, stats.Metric)
		metricTdigestSampleCountKey := getMetricTdigestMetricSamplesKey(stats.OrgID, stats.TableName, stats.Metric)
		pipe.Incr(ctx, metricTdigestSampleCountKey)
		pipe.Expire(ctx, metricTdigestSampleCountKey, StatsTimeToLiveAfterLastUsage)
		tdigestResp = pipe.Do(ctx, "PERCENTILE.MERGE", tableMetricKey, stats.LatencyMilliseconds)
		pipe.Expire(ctx, tableMetricKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes

		tblMetricSet := getTableMetricSetKey(stats.OrgID, stats.TableName)
		pipe.SAdd(ctx, tblMetricSet, stats.Metric)
		pipe.Expire(ctx, tblMetricSet, StatsTimeToLiveAfterLastUsage)
	}
	pipe.Expire(ctx, rateLimiterKey, 2*time.Second)
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

func (m *IrisCache) GetNextAdaptiveRoute(ctx context.Context, orgID int, rgName, metricName string, ri iris_models.RouteInfo, meter *iris_usage_meters.PayloadSizeMeter) (*StatTable, error) {
	ts := &StatTable{
		OrgID:     orgID,
		TableName: rgName,
		MemberRankScoreIn: redis.Z{
			Score:  1,
			Member: ri.RoutePath,
		},
		MemberRankScoreOut:            redis.Z{},
		LatencyQuartilePercentageRank: 0,
		LatencyMilliseconds:           0,
		Metric:                        metricName,
		MetricLatencyMedian:           0,
		MetricLatencyTail:             0,
		MetricSampleCount:             0,
		Meter:                         meter,
	}
	err := m.GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing(ctx, ts)
	if err != nil {
		log.Err(err).Msgf("GetNextAdaptiveRoute")
		return nil, err
	}
	return ts, nil
}

const RateLimitTTL = 2 * time.Second

func (m *IrisCache) SetLatestAdaptiveEndpointPriorityScore(ctx context.Context, stats *StatTable) error {
	if stats == nil {
		return fmt.Errorf("stats is nil")
	}
	if stats.TableName == "" {
		return fmt.Errorf("stats.TableName is empty")
	}

	//endpointPriorityScoreKey := createAdaptiveEndpointPriorityScoreKey(stats.OrgID, stats.TableName)
	//if stats.LatencyQuartilePercentageRank <= 0.0 {
	//	if stats.LatencyMilliseconds > int64(stats.MetricLatencyTail) {
	//		stats.LatencyQuartilePercentageRank = 1.0
	//	} else if stats.LatencyMilliseconds < int64(stats.MetricLatencyMedian) {
	//		stats.LatencyQuartilePercentageRank = 0.5
	//	} else {
	//		stats.LatencyQuartilePercentageRank = 0.75
	//	}
	//}

	//log.Info().Float64(" stats.MetricLatencyMedian", stats.MetricLatencyMedian).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")
	//log.Info().Float64(" stats.MetricLatencyTail", stats.MetricLatencyTail).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")
	//log.Info().Int64(" stats.LatencyMilliseconds", stats.LatencyMilliseconds).Msgf("SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage: latency metrics")
	//stats.ScaleFactor = 0.6
	//rate := stats.LatencyQuartilePercentageRank + stats.ScaleFactor
	// essentially this just multiplies the score by the priority rate growth
	//scoreAdjustmentMemberOut := rate * stats.MemberRankScoreOut.Score
	//stats.MemberRankScoreOut.Score = scoreAdjustmentMemberOut
	pipe := m.Writer.TxPipeline()

	//rateLimiterKey := getOrgRateLimitKey(stats.OrgID)
	//pipe.Expire(ctx, rateLimiterKey, RateLimitTTL)

	//if stats.MetricSampleCount >= MinSamplesBeforeAdaptiveScoring {
	//	pipe.ZAdd(ctx, endpointPriorityScoreKey, stats.MemberRankScoreOut)
	//}
	//if stats.MemberRankScoreIn.Score > 1 {
	//	stats.MemberRankScoreIn.Score *= DecayConstant
	//	pipe.ZAdd(ctx, endpointPriorityScoreKey, stats.MemberRankScoreIn)
	//}

	var tdigestResp *redis.Cmd
	if stats.Metric != "" && stats.LatencyMilliseconds > 0 {
		tableMetricKey := getMetricTdigestKey(stats.OrgID, stats.TableName, stats.Metric)
		metricTdigestSampleCountKey := getMetricTdigestMetricSamplesKey(stats.OrgID, stats.TableName, stats.Metric)
		pipe.Incr(ctx, metricTdigestSampleCountKey)
		pipe.Expire(ctx, metricTdigestSampleCountKey, StatsTimeToLiveAfterLastUsage)
		tdigestResp = pipe.Do(ctx, "PERCENTILE.MERGE", tableMetricKey, stats.LatencyMilliseconds)
		pipe.Expire(ctx, tableMetricKey, StatsTimeToLiveAfterLastUsage) // Set the TTL to 15 minutes

		tblMetricSet := getTableMetricSetKey(stats.OrgID, stats.TableName)
		pipe.SAdd(ctx, tblMetricSet, stats.Metric)
		pipe.Expire(ctx, tblMetricSet, StatsTimeToLiveAfterLastUsage)
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
