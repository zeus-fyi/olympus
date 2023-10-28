package kronos_helix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"go.temporal.io/sdk/client"
)

func (t *KronosWorkerTestSuite) TestMonitorWorkflowStep() {
	ojs, jerr := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, "olympus")
	t.Require().Nil(jerr)
	count := 0
	for _, ojob := range ojs {
		ins := Instructions{}
		err := json.Unmarshal([]byte(ojob.Instructions), &ins)
		t.Require().Nil(err)

		if ins.Type == monitoring {
			count++
		}
	}
	t.Assert().NotZero(count)
}

func (t *KronosWorkerTestSuite) TestMonitors() {
	t.TestHestiaMonitor()
	t.TestZeusMonitor()
	t.TestIrisMonitor()
	t.TestZeusCloudMonitor()
}

func (t *KronosWorkerTestSuite) TestZeusCloudMonitor() (*artemis_orchestrations.OrchestrationJob, Instructions) {
	groupName := "ZeusCloud"
	endpoint := ZeusCloudHealthEndpoint
	pollInterval := time.Second * 30

	return t.testCreateNewMonitorOrchestrationJob(groupName, endpoint, pollInterval)
}

func (t *KronosWorkerTestSuite) TestHestiaMonitor() (*artemis_orchestrations.OrchestrationJob, Instructions) {
	groupName := "Hestia"
	endpoint := HestiaHealthEndpoint
	pollInterval := time.Second * 30

	return t.testCreateNewMonitorOrchestrationJob(groupName, endpoint, pollInterval)
}

func (t *KronosWorkerTestSuite) TestZeusMonitor() (*artemis_orchestrations.OrchestrationJob, Instructions) {
	groupName := "Zeus"
	endpoint := ZeusHealthEndpoint
	pollInterval := time.Second * 30

	return t.testCreateNewMonitorOrchestrationJob(groupName, endpoint, pollInterval)
}

func (t *KronosWorkerTestSuite) TestIrisMonitor() (*artemis_orchestrations.OrchestrationJob, Instructions) {
	groupName := "Iris"
	endpoint := IrisHealthEndpoint
	pollInterval := time.Second * 30

	return t.testCreateNewMonitorOrchestrationJob(groupName, endpoint, pollInterval)
}

func (t *KronosWorkerTestSuite) testCreateNewMonitorOrchestrationJob(groupName, endpoint string, pollInterval time.Duration) (*artemis_orchestrations.OrchestrationJob, Instructions) {
	instType := "HealthMonitor"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		Monitors:  CreateNewMonitorInstructions(groupName, endpoint, pollInterval, 12),
	}
	b, err := json.Marshal(inst)
	t.Require().Nil(err)
	groupName = olympus
	instType = monitoring
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
	t.Require().Equal(groupName, oj.GroupName)
	t.Require().Equal(instType, oj.Type)
	err = oj.UpsertOrchestrationWithInstructions(ctx)
	t.Require().Nil(err)
	t.Assert().NotZero(oj.OrchestrationID)

	return &oj, inst
}

func (t *KronosWorkerTestSuite) TestCreateNewMonitorInstructions() MonitorInstructions {
	pollInterval := time.Second * 30

	serviceName := "iris"
	endpoint := "https://iris.zeus.fyi/health"
	monitorInstructions := CreateNewMonitorInstructions(serviceName, endpoint, pollInterval, 12)
	t.Assert().Equal(serviceName, monitorInstructions.ServiceName)
	t.Assert().Equal(endpoint, monitorInstructions.Endpoint)
	t.Assert().Equal(pollInterval, monitorInstructions.PollInterval)
	t.Assert().Equal(10, monitorInstructions.AlertFailureThreshold)
	return monitorInstructions
}

func (t *KronosWorkerTestSuite) TestMonitorPoll() {
	mi := Instructions{
		GroupName: olympus,
		Type:      monitoring,
		Monitors:  t.TestCreateNewMonitorInstructions(),
	}
	failureCount := 0
	pollCycles := t.TestPollCyclesCount()
	for i := 0; i < pollCycles; i++ {
		resp, err := http.Get(mi.Monitors.Endpoint)
		if err != nil {
			failureCount++
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			failureCount++
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			failureCount = 0
		}
		if failureCount >= mi.Monitors.AlertFailureThreshold {
			fmt.Println("failureCount", failureCount)
		}
		fmt.Println(failureCount)
		time.Sleep(mi.Monitors.PollInterval)
	}
}

func (t *KronosWorkerTestSuite) TestPollCyclesCount() int {
	resetTime := time.Minute * 10
	pollInterval := time.Second * 30
	cycles := CalculatePollCycles(resetTime, pollInterval)
	t.Assert().Equal(20, cycles)
	return cycles
}

func (t *KronosWorkerTestSuite) TestMonitorEndpoint() {
	resp, err := http.Get(ZeusCloudHealthEndpoint)
	t.Require().Nil(err)

	t.Assert().Equal(200, resp.StatusCode)
}

func (t *KronosWorkerTestSuite) TestKronosMonitorWorkflow() {
	ta := t.Tc.DevTemporalAuth
	//ns := "kronos.ngb72"
	//hp := "kronos.ngb72.tmprl.cloud:7233"
	//ta.Namespace = ns
	//ta.HostPort = hp
	InitKronosHelixWorker(ctx, ta)
	cKronos := KronosServiceWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	KronosServiceWorker.Worker.RegisterWorker(cKronos)
	err := KronosServiceWorker.Worker.Start()
	t.Require().Nil(err)

	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: KronosHelixTaskQueue,
	}
	txWf := NewKronosWorkflow()
	wf := txWf.Monitor

	oj, inst := t.TestIrisMonitor()
	cycles := 10
	_, err = cKronos.ExecuteWorkflow(ctx, workflowOptions, wf, oj, inst, cycles)
	t.Require().NoError(err)
}
