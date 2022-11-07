package workers

import temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"

type TopologyWorker struct {
	temporal_base.Worker
}

func NewTopologyWorker(authCfg temporal_base.TemporalAuth) (TopologyWorker, error) {
	tc, err := temporal_base.NewTemporalClient(authCfg)
	if err != nil {
		return TopologyWorker{}, err
	}
	w := temporal_base.NewWorker(tc, "TopologyTaskQueue")
	tw := TopologyWorker{w}
	return tw, err
}
