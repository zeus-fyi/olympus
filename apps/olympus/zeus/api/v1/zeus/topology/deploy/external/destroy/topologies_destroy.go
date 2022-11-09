package destroy_deploy_request

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/helpers"
)

type TopologyDestroyDeployRequest struct {
	TopologyID    int    `db:"topology_id" json:"topologyID"`
	CloudProvider string `db:"cloud_provider" json:"cloudProvider"`
	Region        string `db:"region" json:"region"`
	Context       string `db:"context" json:"context"`
	Namespace     string `db:"namespace" json:"namespace"`
	Env           string `db:"env" json:"env"`
}

func (t *TopologyDestroyDeployRequest) DestroyDeployedTopology(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = t.TopologyID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	// from auth lookup
	bearer := c.Get("bearer")
	knsDestroyDeploy := kns.NewKns()
	knsDestroyDeploy.TopologiesKns = autogen_bases.TopologiesKns{
		TopologyID:    t.TopologyID,
		CloudProvider: t.CloudProvider,
		Region:        t.Region,
		Context:       t.Context,
		Namespace:     t.Namespace,
		Env:           t.Env,
	}
	tar := helpers.PackageCommonTopologyRequest(knsDestroyDeploy, bearer.(string), ou, tr.GetNativeK8s())
	err = topology_worker.Worker.ExecuteDestroyDeploy(ctx, tar)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, ExecuteWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := topology_deployment_status.NewTopologyStatus()
	resp.TopologyID = t.TopologyID
	resp.TopologyStatus = "Pending"
	resp.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusAccepted, resp)
}
