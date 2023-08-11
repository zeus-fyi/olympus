package kronos_helix

import (
	"context"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
)

type KronosWorker struct {
	temporal_base.Worker
}

var (
	KronosServiceWorker KronosWorker
)

const (
	KronosHelixTaskQueue = "KronosHelixTaskQueue"
)

func (k *KronosWorker) ExecuteKronosWorkflow(ctx context.Context) error {
	return nil
}
