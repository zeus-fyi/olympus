package topology_worker

import (
	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
)

func InitTopologyWorker(authCfg temporal_base.TemporalAuth) error {
	var err error
	Worker, err = NewTopologyWorker(authCfg)
	log.Err(err).Msg("InitTopologyWorker failed")
	return err
}
