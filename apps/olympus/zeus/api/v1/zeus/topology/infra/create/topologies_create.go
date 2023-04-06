package create_infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/infra"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

type TopologyCreateRequest struct {
	TopologyName     string `json:"topologyName"`
	ChartName        string `json:"chartName"`
	ChartDescription string `json:"chartDescription,omitempty"`
	Version          string `json:"version"`

	ClusterClassName  string `json:"clusterClassName,omitempty"`
	ComponentBaseName string `json:"componentBaseName,omitempty"`
	SkeletonBaseName  string `json:"skeletonBaseName,omitempty"`
	Tag               string `json:"tag,omitempty"`
}

type TopologyCreateResponse struct {
	TopologyID       int    `json:"topologyID"`
	SkeletonBaseName string `json:"skeletonBaseName,omitempty"`
}

type TopologyCreateRequestFromUI struct {
	zeus_templates.Cluster `json:"cluster"`
}

var ts chronos.Chronos

func (t *TopologyCreateRequestFromUI) CreateTopologyFromUI(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	pcg, err := zeus_templates.GenerateSkeletonBaseChartsPreview(ctx, t.Cluster)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error generating skeleton base charts")
		return c.JSON(http.StatusBadRequest, err)
	}
	gcd, err := zeus_templates.GenerateClusterFromUI(ctx, t.Cluster)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error generating skeleton base charts")
		return c.JSON(http.StatusBadRequest, err)
	}
	bearer := c.Get("bearer").(string)
	zc := zeus_client.NewDefaultZeusClient(bearer)
	// TODO replace this, so that it is wrapped in the tx
	err = gcd.CreateClusterClassDefinitions(ctx, zc)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating cluster class definitions")
		return c.JSON(http.StatusBadRequest, err)
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	defer tx.Rollback(ctx)

	for componentBaseName, component := range pcg.ComponentBases {
		for skeletonBaseName, skeleton := range component {
			nk := chart_workload.TopologyBaseInfraWorkload{}

			if skeleton.Deployment != nil {
				b, berr := json.Marshal(skeleton.Deployment)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			if skeleton.StatefulSet != nil {
				b, berr := json.Marshal(skeleton.StatefulSet)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			if nk.StatefulSet != nil && nk.Deployment != nil {
				err = errors.New("cannot include both a stateful set and deployment, must only choose one per topology infra chart components")
				return c.JSON(http.StatusBadRequest, err)
			}

			if skeleton.Service != nil {
				b, berr := json.Marshal(skeleton.Service)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			if skeleton.Ingress != nil {
				b, berr := json.Marshal(skeleton.Ingress)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			if skeleton.ConfigMap != nil {
				b, berr := json.Marshal(skeleton.ConfigMap)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			if skeleton.ServiceMonitor != nil {
				b, berr := json.Marshal(skeleton.ServiceMonitor)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, nil)
				}
			}

			cw, cerr := nk.CreateChartWorkloadFromTopologyBaseInfraWorkload()
			if cerr != nil {
				log.Err(cerr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: CreateTopology, CreateChartWorkloadFromTopologyBaseInfraWorkload")
				return c.JSON(http.StatusBadRequest, nil)
			}
			inf := create_infra.NewCreateInfrastructure()
			inf.ChartWorkload = cw
			inf.ClusterClassName = t.ClusterName
			inf.ComponentBaseName = componentBaseName
			inf.SkeletonBaseName = skeletonBaseName

			inf.OrgID = ou.OrgID
			inf.UserID = ou.UserID
			inf.Name = skeletonBaseName
			inf.Chart.ChartName = skeletonBaseName
			inf.ChartVersion = fmt.Sprintf("%d", ts.UnixTimeStampNow())
			inf.Tag = "latest"
			tx, err = inf.InsertInfraBaseTx(ctx, tx)
			if err != nil {
				pgErr := err.(*pgconn.PgError)
				switch {
				case strings.Contains(pgErr.Error(), "chart_package_unique"):
					err = errors.New("chart name and version already exists")
					return c.JSON(http.StatusBadRequest, err)
				default:
					log.Err(err).Msg("CreateTopologyFromUI: CreateTopology, InsertInfraBase")
					err = errors.New("unable to add chart, verify it is a valid kubernetes workload that's supported")
				}
				log.Err(err).Interface("orgUser", ou).Msg("TopologyActionCreateRequest: CreateTopology, InsertInfraBase")
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error committing transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, pcg)
}

func (t *TopologyCreateRequest) CreateTopology(c echo.Context) error {
	nk, err := zeus.DecompressUserInfraWorkload(c)
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: CreateTopology, DecompressUserInfraWorkload")
		return c.JSON(http.StatusBadRequest, nil)
	}
	cw, err := nk.CreateChartWorkloadFromTopologyBaseInfraWorkload()
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: CreateTopology, CreateChartWorkloadFromTopologyBaseInfraWorkload")
		return c.JSON(http.StatusBadRequest, nil)
	}
	if nk.StatefulSet != nil && nk.Deployment != nil {
		err = errors.New("cannot include both a stateful set and deployment, must only choose one per topology infra chart components")
		return c.JSON(http.StatusBadRequest, err)
	}
	inf := create_infra.NewCreateInfrastructure()
	ctx := context.Background()
	inf.ChartWorkload = cw

	topologyName := c.FormValue("topologyName")
	inf.Name = topologyName

	chartName := c.FormValue("chartName")
	inf.Chart.ChartName = chartName

	chartDescription := c.FormValue("chartDescription")
	inf.Chart.ChartDescription.String = chartDescription

	version := c.FormValue("version")
	inf.ChartVersion = version

	inf.Tag = c.FormValue("tag")

	// from auth lookup
	ou := c.Get("orgUser").(org_users.OrgUser)
	inf.OrgID = ou.OrgID
	inf.UserID = ou.UserID
	inf.ChartVersion = version

	inf.SkeletonBaseName = c.FormValue("skeletonBaseName")
	inf.ComponentBaseName = c.FormValue("componentBaseName")
	inf.ClusterClassName = c.FormValue("clusterClassName")

	err = inf.InsertInfraBase(ctx)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch {
		case strings.Contains(pgErr.Error(), "chart_package_unique"):
			err = errors.New("chart name and version already exists")
			return c.JSON(http.StatusBadRequest, err)
		default:
			log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: CreateTopology, InsertInfraBase")
			err = errors.New("unable to add chart, verify it is a valid kubernetes workload that's supported")
		}
		log.Err(err).Interface("orgUser", ou).Msg("TopologyActionCreateRequest: CreateTopology, InsertInfraBase")
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := TopologyCreateResponse{
		TopologyID:       inf.TopologyID,
		SkeletonBaseName: inf.SkeletonBaseName,
	}
	return c.JSON(http.StatusOK, resp)
}
