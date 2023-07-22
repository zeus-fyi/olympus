package redis_mev

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
)

const (
	P2PMainnetNodesKey = "p2pMainnetNodes"
)

func (m *MevCache) AddOrUpdateNodeCache(ctx context.Context, nodes artemis_mev_models.P2PNodes, ttl time.Duration) error {
	key := P2PMainnetNodesKey
	// convert P2PNodes to JSON
	nodesJson, err := json.Marshal(nodes)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateNodeCache: Unable to marshal P2PNodes to JSON")
		return err
	}
	statusCmd := m.Set(ctx, key, nodesJson, ttl)
	if statusCmd.Err() != nil {
		log.Err(statusCmd.Err()).Msgf("AddOrUpdateNodeCache: %s", key)
		return statusCmd.Err()
	}
	return nil
}

func (m *MevCache) GetNodeCache() (artemis_mev_models.P2PNodes, error) {
	key := P2PMainnetNodesKey
	// Get JSON string from Redis
	nodesJson, err := m.Get(context.Background(), key).Result()
	if err != nil {
		log.Err(err).Msgf("GetNodeCache: %s", key)
		return nil, err
	}

	// Convert JSON string to P2PNodes
	var nodes artemis_mev_models.P2PNodes
	err = json.Unmarshal([]byte(nodesJson), &nodes)
	if err != nil {
		log.Err(err).Msgf("GetNodeCache: Unable to unmarshal JSON to P2PNodes")
		return nil, err
	}

	return nodes, nil
}
