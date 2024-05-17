package create_or_update_deploy

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	hestia_ovhcloud "github.com/zeus-fyi/olympus/pkg/hestia/ovhcloud"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"k8s.io/apimachinery/pkg/util/validation"
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
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	FreeTrial                    bool                       `json:"freeTrial"`
	MonthlyCost                  float64                    `json:"monthlyCost"`
	Disk                         autogen_bases.Disks        `json:"disk,omitempty"`
	Node                         autogen_bases.Nodes        `json:"nodes"`
	Count                        float64                    `json:"count"`
	NamespaceAlias               string                     `json:"namespaceAlias"`
	Cluster                      zeus_templates.Cluster     `json:"cluster"`
	ResourceRequirements         []DiskResourceRequirements `json:"resourceRequirements"`
	IsPublic                     bool                       `json:"isPublic"`
	AppTaint                     bool                       `json:"appTaint"`
}

func (t *TopologyDeployUIRequest) Validate(ctx context.Context, ou org_users.OrgUser) error {
	if t.CloudProvider == "" {
		return fmt.Errorf("cloudProvider is required")
	}
	if t.Region == "" {
		return fmt.Errorf("region is required")
	}
	if t.Count > 0 && t.Node.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	if t.Disk.DiskSize > 0 && t.Disk.Type == "" {
		return fmt.Errorf("disk type is required")
	}
	if t.Disk.DiskSize > 0 && t.Node.Slug == "" {
		return fmt.Errorf("node slug is required for disk provisioning")
	}
	if t.Disk.DiskSize > 0 && t.Disk.ResourceStrID == "" {
		return fmt.Errorf("resource id is required for disk provisioning")
	}
	if t.Disk.ResourceStrID != "" {
		rid, err := strconv.Atoi(t.Disk.ResourceStrID)
		if err != nil {
			return fmt.Errorf("invalid resource id")
		}
		t.Disk.ResourceID = rid
	}

	if t.Node.ExtCfgStrID != "" && t.Disk.ExtCfgStrID != "" {
		if t.Node.ExtCfgStrID != t.Disk.ExtCfgStrID {
			return fmt.Errorf("node and disk must belong to the same cloud region group")
		}
	}
	kcs, err := authorized_clusters.SelectAuthedAndPublicClusterConfigsByOrgID(ctx, ou)
	if err != nil {
		log.Err(err).Msg("failed to get authorized and public cluster configs")
		return err
	}
	for _, kc := range kcs {
		if kc.CloudProvider == t.CloudProvider && kc.Region == t.Region && kc.ExtConfigStrID == t.Node.ExtCfgStrID {
			t.Context = kc.Context
			t.ClusterCfgStrID = kc.ExtConfigStrID
			t.IsPublic = kc.IsPublic
			return nil
		}
	}
	return fmt.Errorf("node does not belong to the same cloud region group")
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
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(ctx, ou.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		t.FreeTrial = true
	}
	if ou.UserID == 1685378241971196000 || ou.UserID == 7138958574876245565 || ou.OrgID == 1685378241971196000 {
		t.FreeTrial = false
	} else {
		isFreeTrialOngoing, ferr := hestia_compute_resources.DoesOrgHaveOngoingFreeTrial(ctx, ou.OrgID)
		if ferr != nil {
			log.Error().Err(ferr).Msg("failed to check if org has ongoing free trial")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		log.Info().Interface("isFreeTrialOngoing", isFreeTrialOngoing).Interface("ou", ou).Msg("isFreeTrialOngoing")
		if isFreeTrialOngoing {
			log.Error().Err(err).Msg("org has ongoing free trial")
			return c.JSON(http.StatusPreconditionFailed, nil)
		}
		if t.MonthlyCost > 500 {
			log.Error().Err(err).Msg("free trial value exceeds max allowed")
			return c.JSON(http.StatusPreconditionFailed, nil)
		}
	}
	err = t.Validate(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to validate request")
		return c.JSON(http.StatusBadRequest, nil)
	}

	clusterID := uuid.New()
	suffix := strings.Split(clusterID.String(), "-")[0]
	var cr base_deploy_params.ClusterSetupRequest
	alias := fmt.Sprintf("%s-%s", t.NamespaceAlias, suffix)
	clusterNs := clusterID.String()
	if validation.IsDNS1123Label(alias) == nil {
		clusterNs = alias
	}
	switch t.CloudCtxNs.CloudProvider {
	case "do":
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				ClusterCfgStrID: t.ClusterCfgStrID,
				CloudProvider:   t.CloudProvider,
				Context:         t.Context,
				Region:          t.Region,
				Namespace:       clusterNs,
				Alias:           alias,
				Env:             "",
			},
			Nodes:         t.Node,
			NodesQuantity: t.Count,
			Disks: autogen_bases.DisksSlice{
				t.Disk,
			},
			Cluster:  t.Cluster,
			AppTaint: t.AppTaint,
		}
	case "gcp":
		if strings.HasPrefix(t.Cluster.ClusterName, "sui") {
			t.Node.DiskType = "nvme"
		}
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				ClusterCfgStrID: t.ClusterCfgStrID,
				CloudProvider:   t.CloudProvider,
				Region:          t.Region,
				Context:         t.Context,
				Namespace:       clusterNs,
				Alias:           alias,
			},
			Nodes: autogen_bases.Nodes{
				CloudProvider: t.CloudProvider,
				Region:        t.Region,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks: autogen_bases.DisksSlice{
				t.Disk,
			},
			Cluster:  t.Cluster,
			AppTaint: t.AppTaint,
		}
	case "aws":
		switch strings.HasPrefix("i", t.Node.Slug) {
		case true:
			t.Node.DiskType = "nvme"
		}
		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				ClusterCfgStrID: t.ClusterCfgStrID,
				CloudProvider:   t.CloudProvider,
				Region:          t.Region,
				Context:         t.Context,
				Namespace:       clusterNs,
				Alias:           alias,
				Env:             "",
			},
			Nodes: autogen_bases.Nodes{
				CloudProvider: t.CloudProvider,
				Region:        t.Region,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks: autogen_bases.DisksSlice{
				t.Disk,
			},
			Cluster:  t.Cluster,
			AppTaint: t.AppTaint,
		}
	case "ovh":
		ovhContext := t.Context
		namespace := clusterNs
		appTaint := t.AppTaint
		switch ou.UserID {
		case 7138958574876245565:
			if ou.OrgID == 7138983863666903883 {
				ovhContext = hestia_ovhcloud.OvhInternalContext
				switch t.NamespaceAlias {
				case "info-flows-staging":
					appTaint = false
					namespace = "info-flows-staging"
				case "flows":
					appTaint = false
					namespace = "flows"
				case "redis", "redis-master", "redis-replicas":
					namespace = "redis"
				case "redis-cluster":
					namespace = "redis-cluster"
				case "mainnet-staking":
					namespace = "mainnet-staking"
				case "ephemeral-staking":
					namespace = "ephemeral-staking"
				case "goerli-staking":
					namespace = "goerli-staking"
				case "keydb":
					namespace = "keydb"
				case "tyche":
					namespace = "tyche"
					appTaint = false
				case "poseidon":
					namespace = "poseidon"
					appTaint = false
				case "artemis":
					namespace = "artemis"
					appTaint = false
				case "zeus":
					namespace = "zeus"
					appTaint = false
				case "iris":
					namespace = "iris"
					appTaint = false
				case "txFetcher":
					namespace = "tx-fetcher"
					appTaint = true
				case "hestia":
					namespace = "hestia"
					appTaint = false
				case "hera":
					namespace = "hera"
					appTaint = false
				case "aegis":
					namespace = "aegis"
					appTaint = false
				case "hardhat":
					namespace = "hardhat"
					appTaint = false
				case "docs":
					namespace = "docs"
					appTaint = false
				case "olympus":
					namespace = "olympus"
					appTaint = false
				}
			}
		}

		cr = base_deploy_params.ClusterSetupRequest{
			FreeTrial: t.FreeTrial,
			Ou:        ou,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				ClusterCfgStrID: t.ClusterCfgStrID,
				CloudProvider:   t.CloudProvider,
				Region:          t.Region,
				Context:         ovhContext, // hardcoded for now
				Namespace:       namespace,
				Alias:           alias,
				Env:             "",
			},
			Nodes: autogen_bases.Nodes{
				CloudProvider: t.CloudProvider,
				Region:        t.Region,
				ResourceID:    t.Node.ResourceID,
				Slug:          t.Node.Slug,
			},
			NodesQuantity: t.Count,
			Disks: autogen_bases.DisksSlice{
				t.Disk,
			},
			Cluster:  t.Cluster,
			AppTaint: appTaint,
			IsPublic: t.IsPublic,
		}
	}
	if cr.CloudCtxNs.CheckIfEmpty() {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if len(cr.CloudCtxNs.Alias) == 0 {
		cr.CloudCtxNs.Alias = cr.CloudCtxNs.Namespace
	}
	return zeus.ExecuteCreateSetupClusterWorkflow(c, ctx, cr)
}
