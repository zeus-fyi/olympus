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
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	hestia_ovhcloud "github.com/zeus-fyi/olympus/pkg/hestia/ovhcloud"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

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
	FreeTrial            bool                       `json:"freeTrial"`
	MonthlyCost          float64                    `json:"monthlyCost"`
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

	if ou.UserID == 1685378241971196000 {
		t.FreeTrial = false
	} else {
		if t.FreeTrial {
			isFreeTrialOngoing, err := hestia_compute_resources.DoesOrgHaveOngoingFreeTrial(ctx, ou.OrgID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to check if org has ongoing free trial")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			log.Ctx(ctx).Info().Interface("isFreeTrialOngoing", isFreeTrialOngoing).Interface("ou", ou).Msg("isFreeTrialOngoing")
			if isFreeTrialOngoing {
				log.Ctx(ctx).Error().Err(err).Msg("org has ongoing free trial")
				return c.JSON(http.StatusPreconditionFailed, nil)
			}
			if t.MonthlyCost > 500 {
				log.Ctx(ctx).Error().Err(err).Msg("free trial value exceeds max allowed")
				return c.JSON(http.StatusPreconditionFailed, nil)
			}
		}
		if ou.UserID != 7138958574876245565 && !t.FreeTrial {
			isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(ctx, ou.UserID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to check if user has billing method")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			if !isBillingSetup {
				log.Ctx(ctx).Error().Err(err).Msg("user does not have billing method")
				return c.JSON(http.StatusForbidden, nil)
			}
		}
	}

	clusterID := uuid.New()
	suffix := strings.Split(clusterID.String(), "-")[0]
	var cr base_deploy_params.ClusterSetupRequest
	var diskResourceID int
	switch t.CloudProvider {
	case "do":
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "nyc1",
				Context:       "do-nyc1-do-nyc1-zeus-demo", // hardcoded for now
				Namespace:     clusterID.String(),
				Alias:         fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix),
				Env:           "",
			},
			Nodes: autogen_bases.Nodes{
				Region:        t.Region,
				CloudProvider: t.CloudProvider,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks:         autogen_bases.DisksSlice{},
			Cluster:       t.Cluster,
		}
		diskResourceID = 1681408541855876000
	case "gcp":
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "gcp",
				Region:        "us-central1",
				Context:       "gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0", // hardcoded for now
				Namespace:     clusterID.String(),
				Alias:         fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix),
				Env:           "",
			},
			Nodes: autogen_bases.Nodes{
				Region:        t.Region,
				CloudProvider: t.CloudProvider,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks:         autogen_bases.DisksSlice{},
			Cluster:       t.Cluster,
		}
		diskResourceID = 1683165785839881000
	case "aws":
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "aws",
				Region:        "us-west-1",
				Context:       "zeus-us-west-1", // hardcoded for now
				Namespace:     clusterID.String(),
				Alias:         fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix),
				Env:           "",
			},
			Nodes: autogen_bases.Nodes{
				Region:        t.Region,
				CloudProvider: t.CloudProvider,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks:         autogen_bases.DisksSlice{},
			Cluster:       t.Cluster,
		}
		diskResourceID = 1683860918169422000
	case "ovh":
		ovhContext := hestia_ovhcloud.OvhSharedContext
		namespace := clusterID.String()
		switch ou.UserID {
		case 7138958574876245565:
			if ou.OrgID == 7138983863666903883 {
				ovhContext = hestia_ovhcloud.OvhInternalContext
				switch t.NamespaceAlias {
				case "artemis":
					namespace = "artemis"
				case "zeus":
					namespace = "zeus"
				case "iris":
					namespace = "iris"
				case "hestia":
					namespace = "hestia"
				case "hera":
					namespace = "hera"
				case "aegis":
					namespace = "aegis"
				case "hardhat":
					namespace = "hardhat"
				}
			}
		}

		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "ovh",
				Region:        hestia_ovhcloud.OvhRegionUsWestOr1,
				Context:       ovhContext, // hardcoded for now
				Namespace:     namespace,
				Alias:         fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix),
				Env:           "",
			},
			Nodes: autogen_bases.Nodes{
				Region:        t.Region,
				CloudProvider: t.CloudProvider,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks:         autogen_bases.DisksSlice{},
			Cluster:       t.Cluster,
		}
		diskResourceID = 1687637679066833000
	}

	ds := make(autogen_bases.DisksSlice, len(t.ResourceRequirements))
	for i, dr := range t.ResourceRequirements {
		disk := autogen_bases.Disks{
			ResourceID:    diskResourceID,
			Region:        t.Region,
			CloudProvider: t.CloudProvider,
			DiskUnits:     dr.ResourceSumsDisk,
		}
		ds[i] = disk
	}
	cr.Disks = ds
	return zeus.ExecuteCreateSetupClusterWorkflow(c, ctx, cr)
}
