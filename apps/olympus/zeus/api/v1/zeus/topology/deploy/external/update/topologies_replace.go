package deploy_updates

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func ReplaceTopology(c echo.Context, tar zeus_req_types.TopologyDeployRequest) error {
	log.Debug().Msg("TopologyReplaceTopology")
	nk, err := zeus.DecompressUserInfraWorkload(c)
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyReplaceTopology: DecompressUserInfraWorkload")
		return c.JSON(http.StatusBadRequest, nil)
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr, err := zeus.ReadUserTopologyConfig(ctx, tar.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReplaceTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	diffReplacement := zeus.DiffChartUpdate(nk, tr.GetTopologyBaseInfraWorkloadOld())

	tmp := topology_workloads.TopologyBaseInfraWorkload{}
	tmp.CronJob = diffReplacement.CronJob
	tmp.Job = diffReplacement.Job
	tmp.ConfigMap = diffReplacement.ConfigMap
	tmp.Deployment = diffReplacement.Deployment
	tmp.Ingress = diffReplacement.Ingress
	tmp.Service = diffReplacement.Service
	tmp.ServiceMonitor = diffReplacement.ServiceMonitor
	tmp.StatefulSet = diffReplacement.StatefulSet
	return zeus.ExecuteDeployWorkflow(c, ctx, ou, tar, tmp, false, tar.ClusterClassName, tar.ClusterClassName)
}

type DeployClusterUpdateRequestUI struct {
	ClusterClassName string                        `json:"clusterClassName"`
	ClustersDeployed []ClusterTopologyAtCloudCtxNs `json:"clustersDeployed"`
	AppTaint         bool                          `json:"appTaint,omitempty"`
}

type ClusterTopologyAtCloudCtxNs struct {
	TopologyID        int    `json:"topologyID"`
	ClusterName       string `json:"clusterName"`
	ComponentBaseName string `json:"componentBaseName"`
	SkeletonBaseName  string `json:"skeletonBaseName"`
}

func (t *DeployClusterUpdateRequestUI) TopologyUpdateRequestUI(c echo.Context) error {
	log.Debug().Msg("DeployClusterUpdateRequestUI")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	cctxID := c.Request().Header.Get("CloudCtxNsID")
	cID, err := strconv.Atoi(cctxID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	authed, cctx, err := read_topology.IsOrgCloudCtxNsAuthorizedFromID(ctx, ou.OrgID, cID)
	if authed != true {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	existingTopologyIDs, sbOptions := getSkeletonBaseNamesByClusterClassName(t.ClusterClassName, t.ClustersDeployed)
	m := make(map[string]map[string]bool)
	for _, val := range t.ClustersDeployed {
		if _, ok := m[val.ComponentBaseName]; !ok {
			m[val.ComponentBaseName] = make(map[string]bool)
		}
		m[val.ComponentBaseName][val.SkeletonBaseName] = true
	}
	cl, err := read_topology.SelectClusterTopologyFiltered(ctx, ou.OrgID, t.ClusterClassName, sbOptions, m)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DeployClusterTopology: SelectClusterTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	latestClTops := make(map[int]read_topology.ClusterTopologies)
	for _, top := range cl.Topologies {
		latestClTops[top.TopologyID] = top
	}
	var newTopIDs []int
	for key, v := range latestClTops {
		if _, exists := existingTopologyIDs[key]; !exists {
			newTopIDs = append(newTopIDs, key)
			log.Info().Str("clusterName", t.ClusterClassName).Interface("replacing", v).Msg("DeployClusterTopology: replacing")
		}
	}
	if len(newTopIDs) == 0 {
		return c.JSON(http.StatusOK, nil)
	}
	clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterClassName: t.ClusterClassName,
		TopologyIDs:      newTopIDs,
		CloudCtxNs:       cctx,
		OrgUser:          ou,
		AppTaint:         t.AppTaint,
	}
	log.Info().Interface("clDeploy", clDeploy).Msg("TopologyUpdateRequestUI")
	return zeus.ExecuteDeployClusterWorkflow(c, ctx, clDeploy)
}

func getSkeletonBaseNamesByClusterClassName(clusterClassName string, clustersDeployedTopologies []ClusterTopologyAtCloudCtxNs) (map[int]bool, []string) {
	var names []string
	m := make(map[int]bool)
	for _, cluster := range clustersDeployedTopologies {
		if cluster.ClusterName == clusterClassName {
			names = append(names, cluster.SkeletonBaseName)
			m[cluster.TopologyID] = true
		}
	}
	return m, names
}
