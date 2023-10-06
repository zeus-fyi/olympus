package read_infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	create_clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/cluster"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

type PublicAppsMatrixRequest struct {
}

func PublicAppsMatrixRequestHandler(c echo.Context) error {
	request := new(PublicAppsMatrixRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetPublicAppFamily(c)
}

func (a *PublicAppsMatrixRequest) GetPublicAppFamily(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Err(nil).Msg("GetAppFamily: orgUser not found")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ctx := context.Background()
	appList, err := read_topology.SelectPublicMatrixApps(ctx, AppsOrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	go func(al autogen_bases.TopologySystemComponentsSlice) {
		for _, app := range al {
			err = CopyAppFamily(ctx, app, ou)
			if err != nil {
				log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
			}
		}
	}(appList)
	return c.JSON(http.StatusOK, appList)
}

func CopyAppFamily(ctx context.Context, appDetails autogen_bases.TopologySystemComponents, ou org_users.OrgUser) error {
	resp := TopologyUIAppDetailsResponse{
		Cluster: zeus_templates.Cluster{
			ClusterName:     appDetails.TopologySystemComponentName,
			ComponentBases:  make(map[string]zeus_templates.SkeletonBases),
			IngressSettings: zeus_templates.Ingress{},
			IngressPaths:    nil,
		},
		ClusterPreviewWorkloadsOlympus: zeus_templates.ClusterPreviewWorkloadsOlympus{
			ClusterName:    appDetails.TopologySystemComponentName,
			ComponentBases: make(map[string]map[string]any),
		},
		SelectedComponentBaseName: "",
		SelectedSkeletonBaseName:  "",
	}
	selectedApp, err := read_topology.SelectAppTopologyByID(ctx, AppsOrgID, appDetails.TopologySystemComponentID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return err
	}
	pcg, err := zeus_templates.GenerateSkeletonBaseChartsCopy(ctx, &selectedApp)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return err
	}
	rsMinMax := zeus_core.ResourceMinMax{
		Max: zeus_core.ResourceAggregate{},
		Min: zeus_core.ResourceAggregate{},
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return err
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
		return err
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
				return err
			}
			sb.ResourceSums = rs
			sbUI[sbName] = sb
		}
		resp.Cluster.ComponentBases[cbName] = sbUI
		resp.ClusterPreviewWorkloadsOlympus.ComponentBases[cbName] = uiSbs
	}
	return err
}
