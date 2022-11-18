package destroy_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
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
	log.Debug().Msg("DestroyDeployedTopology")
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
	knsDestroyDeploy := kns.NewKns()
	knsDestroyDeploy.TopologiesKns = autogen_bases.TopologiesKns{
		TopologyID:    t.TopologyID,
		CloudProvider: t.CloudProvider,
		Region:        t.Region,
		Context:       t.Context,
		Namespace:     t.Namespace,
		Env:           t.Env,
	}
	// validate context kns
	authed, err := tr.IsOrgCloudCtxNsAuthorized(ctx, knsDestroyDeploy)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return zeus.ExecuteDestroyDeployWorkflow(c, ctx, ou, knsDestroyDeploy, tr.GetNativeK8s())
}
