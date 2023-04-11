package read_infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	create_clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/cluster"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

type AvaxAppsPageRequest struct {
}

func AvaxAppsHandler(c echo.Context) error {
	request := new(AvaxAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetAvaxApp(c)
}

func EthAppsHandler(c echo.Context) error {
	request := new(AvaxAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetEphemeralBeaconsApp(c)
}

const (
	AvaxAppID                = 1680924257606485000
	EphemeralEthBeaconsAppID = 1670997020811171000

	AppsOrgID  = 7138983863666903883
	AppsUserID = 7138958574876245565
)

func (a *AvaxAppsPageRequest) GetEphemeralBeaconsApp(c echo.Context) error {
	return a.GetApp(c, AppsOrgID, EphemeralEthBeaconsAppID)
}
func (a *AvaxAppsPageRequest) GetAvaxApp(c echo.Context) error {
	return a.GetApp(c, AppsOrgID, AvaxAppID)
}
func (a *AvaxAppsPageRequest) GetApp(c echo.Context, AppsOrgID, AppID int) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()

	selectedApp, err := read_topology.SelectAppTopologyByID(ctx, AppsOrgID, AppID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := TopologyUIAppDetailsResponse{
		Cluster: zeus_templates.Cluster{
			ClusterName:     selectedApp.ClusterClassName,
			ComponentBases:  make(map[string]zeus_templates.SkeletonBases),
			IngressSettings: zeus_templates.Ingress{},
			IngressPaths:    nil,
		},
		ClusterPreviewWorkloadsOlympus: zeus_templates.ClusterPreviewWorkloadsOlympus{
			ClusterName:    selectedApp.ClusterClassName,
			ComponentBases: make(map[string]map[string]any),
		},
		SelectedComponentBaseName: "",
		SelectedSkeletonBaseName:  "",
	}
	pcg, _ := zeus_templates.GenerateSkeletonBaseChartsCopy(ctx, &selectedApp)
	rsMinMax := zeus_core.ResourceMinMax{
		Max: zeus_core.ResourceAggregate{},
		Min: zeus_core.ResourceAggregate{},
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	defer tx.Rollback(ctx)
	sys := systems.Systems{TopologySystemComponents: autogen_bases.TopologySystemComponents{
		OrgID:                       ou.OrgID,
		TopologyClassTypeID:         class_types.ClusterClassTypeID,
		TopologySystemComponentName: selectedApp.ClusterClassName,
	}}

	tx, err = create_clusters.InsertCluster(ctx, tx, &sys, pcg, ou)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	for cbName, cb := range pcg.ComponentBases {
		uiSbs := make(map[string]any)
		resp.ClusterPreviewWorkloadsOlympus.ComponentBases[cbName] = make(map[string]any)
		resp.SelectedComponentBaseName = cbName
		sbUI := make(map[string]zeus_templates.SkeletonBase)
		for sbName, nk := range cb {
			sb := zeus_templates.SkeletonBase{}
			resp.SelectedSkeletonBaseName = sbName
			uiSbs[sbName] = nk

			rs := zeus_core.ResourceSums{}
			if nk.StatefulSet != nil {
				sb.AddStatefulSet = true
				if nk.StatefulSet.Spec.Replicas != nil {
					rs.Replicas = fmt.Sprintf("%d", *nk.StatefulSet.Spec.Replicas)
				}
				zeus_core.GetResourceRequirements(ctx, nk.StatefulSet.Spec.Template.Spec, &rs)
				zeus_core.GetDiskRequirements(ctx, nk.StatefulSet.Spec.VolumeClaimTemplates, &rs)
			}
			if nk.Deployment != nil {
				sb.AddDeployment = true
				if nk.Deployment.Spec.Replicas != nil {
					rs.Replicas = fmt.Sprintf("%d", *nk.Deployment.Spec.Replicas)
				}
				zeus_core.GetResourceRequirements(ctx, nk.Deployment.Spec.Template.Spec, &rs)
			}
			if nk.Service != nil {
				sb.AddService = true
			}
			if nk.ConfigMap != nil {
				sb.AddConfigMap = true
			}
			if nk.Ingress != nil {
				sb.AddIngress = true
			}
			if nk.ServiceMonitor != nil {
				sb.AddServiceMonitor = true
			}
			rsMinMax, err = zeus_core.ApplyMinMaxConstraints(rs, rsMinMax)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: ApplyMinMaxConstraints")
				return c.JSON(http.StatusInternalServerError, nil)
			}
			sb.ResourceSums = rs
			sbUI[sbName] = sb
		}
		resp.Cluster.ComponentBases[cbName] = sbUI
		resp.ClusterPreviewWorkloadsOlympus.ComponentBases[cbName] = uiSbs
	}
	nf := hestia_compute_resources.NodeFilter{
		CloudProvider: "do",
		Region:        "nyc1",
		ResourceSums: zeus_core.ResourceSums{
			MemRequests: rsMinMax.Min.MemRequests,
			CpuRequests: rsMinMax.Min.CpuRequests,
		},
	}
	nodes, err := hestia_compute_resources.SelectNodes(ctx, nf)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectNodes")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp.Nodes = nodes
	return c.JSON(http.StatusOK, resp)
}
