package read_infra

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_nodes "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/nodes"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

type TopologyReadPrivateAppsRequest struct {
}

func (t *TopologyReadPrivateAppsRequest) ListPrivateAppsRequest(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	apps, err := read_topology.SelectOrgApps(ctx, ou.OrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, apps)
}

type TopologyUIAppDetailsResponse struct {
	zeus_templates.Cluster                        `json:"cluster"`
	zeus_templates.ClusterPreviewWorkloadsOlympus `json:"clusterPreview"`
	SelectedComponentBaseName                     string                          `json:"selectedComponentBaseName"`
	SelectedSkeletonBaseName                      string                          `json:"selectedSkeletonBaseName"`
	Nodes                                         hestia_autogen_bases.NodesSlice `json:"nodes,omitempty"`
}

func (t *TopologyReadPrivateAppsRequest) GetPrivateAppDetailsRequest(c echo.Context) error {
	appID := c.Param("id")
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	id, err := strconv.Atoi(appID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("GetPrivateAppDetailsRequest: SelectOrgApps")
		return c.JSON(http.StatusBadRequest, nil)
	}
	apps, err := read_topology.SelectAppTopologyByID(ctx, ou.OrgID, id)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := TopologyUIAppDetailsResponse{
		Cluster: zeus_templates.Cluster{
			ClusterName:     apps.ClusterClassName,
			ComponentBases:  make(map[string]zeus_templates.SkeletonBases),
			IngressSettings: zeus_templates.Ingress{},
			IngressPaths:    nil,
		},
		ClusterPreviewWorkloadsOlympus: zeus_templates.ClusterPreviewWorkloadsOlympus{
			ClusterName:    apps.ClusterClassName,
			ComponentBases: make(map[string]map[string]any),
		},
		SelectedComponentBaseName: "",
		SelectedSkeletonBaseName:  "",
	}

	rsMinMax := zeus_core.ResourceMinMax{
		Max: zeus_core.ResourceAggregate{},
		Min: zeus_core.ResourceAggregate{},
	}

	for cbName, cb := range apps.ComponentBases {
		uiSbs := make(map[string]any)
		resp.ClusterPreviewWorkloadsOlympus.ComponentBases[cbName] = make(map[string]any)
		resp.SelectedComponentBaseName = cbName
		sbUI := make(map[string]zeus_templates.SkeletonBase)
		for sbName, sb := range cb.SkeletonBases {
			resp.SelectedSkeletonBaseName = sbName
			tr := read_topology.NewInfraTopologyReader()
			tr.TopologyID = sb.TopologyID
			// from auth lookup
			tr.OrgID = ou.OrgID
			tr.UserID = ou.UserID
			err = tr.SelectTopologyForOrg(ctx)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectTopology")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			nk := tr.GetTopologyBaseInfraWorkload()
			uiSbs[sbName] = nk
			sbTemplate := zeus_templates.SkeletonBase{
				TopologyID: fmt.Sprintf("%d", sb.TopologyID),
			}
			rs := zeus_core.ResourceSums{}
			if nk.StatefulSet != nil {
				sbTemplate.AddStatefulSet = true
				if nk.StatefulSet.Spec.Replicas != nil {
					rs.Replicas = fmt.Sprintf("%d", *nk.StatefulSet.Spec.Replicas)
				}
				zeus_core.GetResourceRequirements(ctx, nk.StatefulSet.Spec.Template.Spec, &rs)
				zeus_core.GetDiskRequirements(ctx, nk.StatefulSet.Spec.VolumeClaimTemplates, &rs)
			}
			if nk.Deployment != nil {
				sbTemplate.AddDeployment = true
				if nk.Deployment.Spec.Replicas != nil {
					rs.Replicas = fmt.Sprintf("%d", *nk.Deployment.Spec.Replicas)
				}
				zeus_core.GetResourceRequirements(ctx, nk.Deployment.Spec.Template.Spec, &rs)
			}
			if nk.Service != nil {
				sbTemplate.AddService = true
			}
			if nk.ConfigMap != nil {
				sbTemplate.AddConfigMap = true
			}
			if nk.Ingress != nil {
				sbTemplate.AddIngress = true
			}
			if nk.ServiceMonitor != nil {
				sbTemplate.AddServiceMonitor = true
			}
			rsMinMax, err = zeus_core.ApplyMinMaxConstraints(rs, rsMinMax)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: ApplyMinMaxConstraints")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			sbTemplate.ResourceSums = rs
			sbUI[sbName] = sbTemplate
		}
		resp.Cluster.ComponentBases[cbName] = sbUI
		resp.ClusterPreviewWorkloadsOlympus.ComponentBases[cbName] = uiSbs
	}
	nf := hestia_nodes.NodeFilter{
		CloudProvider: "do",
		Region:        "nyc1",
		ResourceSums: zeus_core.ResourceSums{
			MemRequests: rsMinMax.Min.MemRequests,
			CpuRequests: rsMinMax.Min.CpuRequests,
		},
	}
	nodes, err := hestia_nodes.SelectNodes(ctx, nf)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectNodes")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp.Nodes = nodes
	return c.JSON(http.StatusOK, resp)
}
