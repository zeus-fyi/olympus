package external_api_workloads

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type TopologyCloudCtxNsQueryRequest struct {
	CloudProvider string `json:"cloudProvider"`
	Region        string `json:"region"`
	Context       string `json:"context"`
	Namespace     string `json:"namespace"`
	Env           string `json:"env"`
}

func (t *TopologyCloudCtxNsQueryRequest) ReadDeployedWorkloads(c echo.Context) error {
	log.Debug().Msg("TopologyCloudCtxNsQueryRequest")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)

	knsDeploy := kns.NewKns()
	knsDeploy.TopologiesKns = autogen_bases.TopologiesKns{
		CloudProvider: t.CloudProvider,
		Region:        t.Region,
		Context:       t.Context,
		Namespace:     t.Namespace,
		Env:           t.Env,
	}
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, knsDeploy)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return nil

}
