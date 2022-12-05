package create_infra

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/infra"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyCreateRequest struct {
	TopologyName     string `json:"topologyName"`
	ChartName        string `json:"chartName"`
	ChartDescription string `json:"chartDescription,omitempty"`
	Version          string `json:"version"`

	SkeletonBaseID int `json:"skeletonBaseID,omitempty"`
}

type TopologyCreateResponse struct {
	TopologyID int `json:"topologyID"`
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

	// from auth lookup
	ou := c.Get("orgUser").(org_users.OrgUser)
	inf.OrgID = ou.OrgID
	inf.UserID = ou.UserID

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
		TopologyID: inf.TopologyID,
	}
	return c.JSON(http.StatusOK, resp)
}
