package iris_redis

import (
	"context"

	"github.com/rs/zerolog/log"
)

type ScaleFactors struct {
	LatencyScaleFactor float64 `json:"latencyScaleFactor,omitempty"`
	ErrorScaleFactor   float64 `json:"errorScaleFactor,omitempty"`
	DecayScaleFactor   float64 `json:"decayScaleFactor,omitempty"`
}

func (m *IrisCache) SetTableLatencyScaleFactor(ctx context.Context, orgID int, rgName string, latSf float64) error {
	latSfKey := createAdaptiveEndpointPriorityScoreLatencyScaleFactorKey(orgID, rgName)
	pipe := m.Writer.TxPipeline()
	pipe.Set(ctx, latSfKey, latSf, 0)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetTableLatencyScaleFactor")
		return err
	}
	return err
}

func (m *IrisCache) SetTableErrorScaleFactor(ctx context.Context, orgID int, rgName string, errSf float64) error {
	errSfKey := createAdaptiveEndpointPriorityScoreErrorScaleFactorKey(orgID, rgName)
	pipe := m.Writer.TxPipeline()
	pipe.Set(ctx, errSfKey, errSf, 0)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetTableErrorScaleFactor")
		return err
	}
	return err
}

func (m *IrisCache) SetTableDecayScaleFactor(ctx context.Context, orgID int, rgName string, decSf float64) error {
	decaySfKey := createAdaptiveEndpointPriorityScoreDecayScaleFactorKey(orgID, rgName)
	pipe := m.Writer.TxPipeline()
	pipe.Set(ctx, decaySfKey, decSf, 0)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetTableDecayScaleFactor")
		return err
	}
	return err
}
