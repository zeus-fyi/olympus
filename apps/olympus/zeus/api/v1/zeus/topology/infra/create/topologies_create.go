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
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/infra"
	create_clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/cluster"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
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
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		err := errors.New("unable to get orgUser from echo context")
		log.Err(err).Msg("CreateTopologyFromUI: CreateTopologyFromUI")
		return c.JSON(http.StatusBadRequest, nil)
	}
	pcg, err := zeus_templates.GenerateSkeletonBaseChartsPreview(ctx, t.Cluster)
	if err != nil {
		log.Err(err).Msg("error generating skeleton base charts")
		return c.JSON(http.StatusBadRequest, err)
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	defer tx.Rollback(ctx)
	sys := systems.Systems{TopologySystemComponents: autogen_bases.TopologySystemComponents{
		OrgID:                       ou.OrgID,
		TopologyClassTypeID:         class_types.ClusterClassTypeID,
		TopologySystemComponentName: t.Cluster.ClusterName,
	}}
	tx, err = create_clusters.InsertCluster(ctx, tx, &sys, pcg, ou)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating transaction")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	defer tx.Rollback(ctx)
	tx, err = apps.Pg.Begin(ctx)
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
			if skeleton.Job != nil {
				b, berr := json.Marshal(skeleton.Job)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, berr)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, err)
				}
			}
			if skeleton.CronJob != nil {
				b, berr := json.Marshal(skeleton.CronJob)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, berr)
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return c.JSON(http.StatusBadRequest, err)
				}
			}
			if nk.Job != nil && nk.CronJob != nil {
				err = errors.New("cannot include both a job and cronjob, must only choose one per topology infra chart components")
				return c.JSON(http.StatusBadRequest, err)
			}
			if (nk.Job != nil || nk.CronJob != nil) && (nk.Deployment != nil || nk.StatefulSet != nil) {
				err = errors.New("cannot include both a job or cronjob with statefulset or deployment, must only choose one class type per infra chart package")
				return c.JSON(http.StatusBadRequest, err)
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
				if strings.Contains(t.IngressSettings.AuthServerURL, "aegis.zeus.fyi") {
					if skeleton.Ingress.Annotations == nil {
						skeleton.Ingress.Annotations = make(map[string]string)
					}
					skeleton.Ingress.Annotations["nginx.ingress.kubernetes.io/auth-url"] = fmt.Sprintf("https://aegis.zeus.fyi/v1/auth/%d", ou.OrgID)
				}
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
	ou := c.Get("orgUser").(org_users.OrgUser)

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
