package artemis_ethereum_transcations

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisEthereumTxBroadcastWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumBroadcastTxActivities
}

const defaultTimeout = 10 * time.Minute

func NewArtemisBroadcastEthereumTxWorkflow() ArtemisEthereumTxBroadcastWorkflow {
	deployWf := ArtemisEthereumTxBroadcastWorkflow{
		Workflow:                             temporal_base.Workflow{},
		ArtemisEthereumBroadcastTxActivities: ArtemisEthereumBroadcastTxActivities{},
	}
	return deployWf
}

func (t *ArtemisEthereumTxBroadcastWorkflow) GetWorkflow() interface{} {
	return t.ArtemisBroadcastEthereumTxWorkflow
}

func (t *ArtemisEthereumTxBroadcastWorkflow) ArtemisBroadcastEthereumTxWorkflow(ctx workflow.Context, params interface{}) error {
	//log := workflow.GetLogger(ctx)

	return nil
}
