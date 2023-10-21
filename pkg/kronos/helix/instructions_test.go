package kronos_helix

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type KronosInstructionsTestSuite struct {
	test_suites_base.TestSuite
}

func (t *KronosInstructionsTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
}

// HestiaPlatformServiceWorkflows
// IrisRoutingServiceRequestWorkflow, IrisDeleteOrgGroupRoutingTableWorkflow, IrisDeleteOrgRoutesWorkflow, IrisRemoveAllOrgRoutesFromCacheWorkflow

// You can change any params for this, it is a template of the other test meant for creating alerts
func (t *KronosWorkerTestSuite) TestInsertAlertOrchestratorsScratchPad() {
	groupName := "ClusterSetupWorkflows"
	instType := "DeployClusterSetupWorkflow"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		Alerts: AlertInstructions{
			Severity:  apollo_pagerduty.CRITICAL,
			Source:    TemporalAlerts,
			Component: "This is a cluster setup workflow for the Zeus subsystem",
			Message:   "A Zeus services workflow is stuck",
		},
		Trigger: TriggerInstructions{
			AlertAfterTime:              time.Minute * 20,
			ResetAlertAfterTimeDuration: time.Minute * 20,
		},
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)
	groupName = olympus
	instType = "alerts"
	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         groupName,
			Type:              instType,
			Instructions:      string(b),
			OrchestrationName: orchName,
		},
	}
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)
}

func (t *KronosWorkerTestSuite) TestInsertAlertOrchestrators() {
	groupName := "HestiaPlatformServiceWorkflows"
	instType := "IrisRemoveAllOrgRoutesFromCacheWorkflow"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		Alerts: AlertInstructions{
			Severity:  apollo_pagerduty.CRITICAL,
			Source:    TemporalAlerts,
			Component: "This is a workflow component",
			Message:   "A QuickNode services workflow is stuck",
		},
		Trigger: TriggerInstructions{
			AlertAfterTime:              time.Minute * 10,
			ResetAlertAfterTimeDuration: time.Minute * 10,
		},
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)
	groupName = olympus
	instType = "alerts"
	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         groupName,
			Type:              instType,
			Instructions:      string(b),
			OrchestrationName: orchName,
		},
	}
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)
}

func (t *KronosInstructionsTestSuite) TestAlertPatternWf() {
	groupName := "testGroup"
	instType := "testType"
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		Alerts: AlertInstructions{
			Severity:  apollo_pagerduty.CRITICAL,
			Message:   "test message",
			Source:    "testSource",
			Component: "testWf",
		},
		Trigger: TriggerInstructions{
			AlertAfterTime:              time.Second * 10,
			ResetAlertAfterTimeDuration: time.Hour,
		},
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)

	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         groupName,
			Type:              instType,
			OrchestrationName: "testKronosMaster-" + uuid.New().String(),
		},
	}
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)

	groupName = olympus
	instType = "alerts"
	oj = artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             t.Tc.ProductionLocalTemporalOrgID,
			Active:            true,
			GroupName:         groupName,
			Type:              instType,
			Instructions:      string(b),
			OrchestrationName: "testKronosAlertTriggerWf-" + uuid.New().String(),
		},
	}

	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)

	selectedOj, err := artemis_orchestrations.SelectOrchestrationByName(ctx, oj.OrgID, oj.OrchestrationName)
	t.Require().Nil(err)

	ka := NewKronosActivities()
	decodedInst, err := ka.GetInstructionsFromJob(ctx, selectedOj)
	t.Require().Nil(err)
	t.Assert().Equal(inst.Type, decodedInst.Type)
	t.Assert().Equal(inst.GroupName, decodedInst.GroupName)
	t.Assert().Equal(inst.Trigger.AlertAfterTime, decodedInst.Trigger.AlertAfterTime)
	t.Assert().Equal(inst.Trigger.ResetAlertAfterTimeDuration, decodedInst.Trigger.ResetAlertAfterTimeDuration)
}

func TestKronosInstructionsTestSuite(t *testing.T) {
	suite.Run(t, new(KronosInstructionsTestSuite))
}
