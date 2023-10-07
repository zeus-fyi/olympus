package read_infra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
)

type PublicAppsPageRequest struct {
}

func MicroserviceAppsHandler(c echo.Context) error {
	request := new(PublicAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetMicroserviceApp(c)
}

func AvaxAppsHandler(c echo.Context) error {
	request := new(PublicAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetAvaxApp(c)
}

func EthAppsHandler(c echo.Context) error {
	request := new(PublicAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetEphemeralBeaconsApp(c)
}

func SuiAppsHandler(c echo.Context) error {
	request := new(PublicAppsPageRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetSuiApp(c)
}

//func AppsByNameHandler(c echo.Context) error {
//	request := new(PublicAppsPageRequest)
//	if err := c.Bind(request); err != nil {
//		return err
//	}
//	return request.GetAppByName(c)
//}

const (
	AvaxAppID                = 1680924257606485000
	EphemeralEthBeaconsAppID = 1670997020811171000
	MicroserviceAppID        = 1681932523630136000
	SuiAppID                 = 1694727626052689000

	AppsOrgID  = 7138983863666903883
	AppsUserID = 7138958574876245565
)

func (a *PublicAppsPageRequest) GetMicroserviceApp(c echo.Context) error {
	return a.GetAppByID(c, MicroserviceAppID)
}
func (a *PublicAppsPageRequest) GetEphemeralBeaconsApp(c echo.Context) error {
	return a.GetAppByID(c, EphemeralEthBeaconsAppID)
}
func (a *PublicAppsPageRequest) GetAvaxApp(c echo.Context) error {
	return a.GetAppByID(c, AvaxAppID)
}
func (a *PublicAppsPageRequest) GetSuiApp(c echo.Context) error {
	return a.GetAppByID(c, SuiAppID)
}

func (a *PublicAppsPageRequest) GetAppByName(c echo.Context, appName string) error {
	ctx := context.Background()
	if !strings.HasPrefix(appName, "sui-") || len(appName) == 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	token, ok := c.Get("bearer").(string)
	if ok {
		err := CopySuiApp(ctx, appName, token)
		if err != nil {
			log.Err(err).Msg("ListPrivateAppsRequest: CopySuiApp")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	selectedApp, err := read_topology.SelectAppTopologyByName(ctx, AppsOrgID, appName)
	if err != nil {
		log.Err(err).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return a.GetApp(c, selectedApp)
}

func (a *PublicAppsPageRequest) GetAppByID(c echo.Context, appID int) error {
	ctx := context.Background()
	selectedApp, err := read_topology.SelectAppTopologyByID(ctx, AppsOrgID, appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return a.GetApp(c, selectedApp)
}

func (a *PublicAppsPageRequest) GetApp(c echo.Context, selectedApp zeus_cluster_config_drivers.ClusterDefinition) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Err(fmt.Errorf("orgUser not found")).Msg("ListPrivateAppsRequest: Get")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ctx := context.Background()
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
	pcg, err := zeus_templates.GenerateSkeletonBaseChartsCopy(ctx, &selectedApp)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
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
				zeus_core.GetBlockStorageDiskRequirements(ctx, nk.StatefulSet.Spec.VolumeClaimTemplates, &rs)
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

	cp := "do"
	region := "nyc1"
	diskType := "ssd"
	switch {
	case strings.Contains(selectedApp.ClusterClassName, "-aws"):
		cp = "aws"
		region = "us-west-1"
		diskType = setNvmeType(selectedApp.ClusterClassName)
	case strings.Contains(selectedApp.ClusterClassName, "-do"):
		cp = "do"
		region = "nyc1"
		diskType = setNvmeType(selectedApp.ClusterClassName)
	case strings.Contains(selectedApp.ClusterClassName, "-gcp"):
		cp = "gcp"
		region = "us-central1"
		diskType = setNvmeType(selectedApp.ClusterClassName)
	case strings.Contains(selectedApp.ClusterClassName, "-ovh"):
		cp = "ovh"
		region = "us-west-or-1"
		diskType = setNvmeType(selectedApp.ClusterClassName)
	default:
		cp = "ovh"
		region = "us-west-or-1"
	}

	nf := hestia_compute_resources.NodeFilter{
		CloudProvider: cp,
		Region:        region,
		ResourceSums: zeus_core.ResourceSums{
			MemRequests: rsMinMax.Min.MemRequests,
			CpuRequests: rsMinMax.Min.CpuRequests,
		},
		DiskType: diskType,
	}
	nodes, err := hestia_compute_resources.SelectNodes(ctx, nf)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectNodes")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if strings.Contains(selectedApp.ClusterClassName, "sui-") {
		for _, node := range nodes {
			switch {
			case strings.Contains(selectedApp.ClusterClassName, "-gcp"):
				node.Disk = 6000
			}
		}
	}
	resp.Nodes = nodes
	return c.JSON(http.StatusOK, resp)
}

func setNvmeType(appName string) string {
	if strings.Contains(appName, "sui-") {
		return "nvme"
	}
	return "ssd"
}
