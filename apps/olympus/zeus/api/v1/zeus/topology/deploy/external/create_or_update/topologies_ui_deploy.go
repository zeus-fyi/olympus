package create_or_update_deploy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

//const payload = {
//"node": node,
//"count": count,
//"namespaceAlias": namespaceAlias,
//"cluster": cluster,
//"resourceRequirements": resourceRequirements,
//}

type DiskResourceRequirements struct {
	ResourceID           string  `json:"resourceID"`
	ComponentBaseName    string  `json:"componentBaseName"`
	SkeletonBaseName     string  `json:"skeletonBaseName"`
	ResourceSumsDisk     string  `json:"resourceSumsDisk"`
	Replicas             string  `json:"replicas"`
	BlockStorageCostUnit float64 `json:"blockStorageCostUnit"`
}

type TopologyDeployUIRequest struct {
	zeus_common_types.CloudCtxNs
	CloudProvider        string                     `json:"cloudProvider"`
	Region               string                     `json:"region"`
	Node                 autogen_bases.Nodes        `json:"nodes"`
	Count                float64                    `json:"count"`
	NamespaceAlias       string                     `json:"namespaceAlias"`
	Cluster              zeus_templates.Cluster     `json:"cluster"`
	ResourceRequirements []DiskResourceRequirements `json:"resourceRequirements"`
}

func SetupClusterTopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyDeployUIRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeploySetupClusterTopology(c)
}

func (t *TopologyDeployUIRequest) DeploySetupClusterTopology(c echo.Context) error {
	log.Debug().Msg("DeploySetupClusterTopology")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(ctx, ou.UserID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		log.Ctx(ctx).Error().Err(err).Msg("user does not have billing method")
		return c.JSON(http.StatusForbidden, nil)
	}
	clusterID := uuid.New()
	suffix := strings.Split(clusterID.String(), "-")[0]
	cr := base_deploy_params.ClusterSetupRequest{
		Ou: ou,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "nyc1",
			Context:       "do-nyc1-do-nyc1-zeus-demo",
			Namespace:     fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix),
			Env:           "",
		},
		ClusterID: clusterID,
		Nodes: autogen_bases.Nodes{
			Region:        t.Region,
			CloudProvider: t.CloudProvider,
			ResourceID:    t.Node.ResourceID,
		},
		NodesQuantity: t.Count,
		Disks:         autogen_bases.DisksSlice{},
	}

	ds := make(autogen_bases.DisksSlice, len(t.ResourceRequirements))
	for i, dr := range t.ResourceRequirements {
		disk := autogen_bases.Disks{
			ResourceID:    1680989153453105000, // hard coded for now
			Region:        t.Region,
			CloudProvider: t.CloudProvider,
			DiskUnits:     dr.ResourceSumsDisk,
		}
		ds[i] = disk
	}
	cr.Disks = ds
	return zeus.ExecuteCreateSetupClusterWorkflow(c, ctx, cr)
}
