package topology_worker

import (
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
)

var Worker TopologyWorker

type TopologyWorker struct {
	temporal_base.Worker
}
