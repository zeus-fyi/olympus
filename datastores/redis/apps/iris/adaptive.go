package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func (m *IrisCache) SetMetricLatencyTDigest(ctx context.Context, orgID int, tableName, metricName string, latency float64) error {
	metricTdigestKey := fmt.Sprintf("%d:%s:%s", orgID, tableName, metricName)

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	pipe.Set(ctx, metricTdigestKey, latency, time.Minute*60)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("SetMetricLatencyTDigest")
		return err
	}
	return nil
}
