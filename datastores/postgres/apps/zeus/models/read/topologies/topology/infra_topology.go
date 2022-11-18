package read_topology

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology
	read_charts.Chart
}

func NewInfraTopologyReader() InfraBaseTopology {
	bt := infra.NewInfrastructureBaseTopology()
	cr := read_charts.NewChartReader()
	rt := InfraBaseTopology{bt, cr}
	return rt
}

func NewInfraTopologyReaderWithOrgUser(ou org_users.OrgUser) InfraBaseTopology {
	bt := infra.NewInfrastructureBaseTopologyWithOrgUser(ou)
	cr := read_charts.NewChartReader()
	rt := InfraBaseTopology{bt, cr}
	return rt
}

const Sn = "ReadInfraBaseTopology"

func (t *InfraBaseTopology) SelectInfraTopologyQuery() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	q.CTEQuery.Params = append(q.CTEQuery.Params, t.TopologyID, t.OrgID, t.UserID)
	q.RawQuery = read_charts.FetchChartQuery(q)
	return q
}

func (t *InfraBaseTopology) SelectTopology(ctx context.Context) error {
	q := t.SelectInfraTopologyQuery()

	log.Debug().Interface("SelectTopologyQuery", q.LogHeader(Sn))
	err := t.SelectSingleChartsResources(ctx, q)
	return err
}

func getIsOrgCloudCtxNsAuthorizedQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("IsOrgCloudCtxNsAuthorized", "topologies_org_cloud_ctx_ns", "where", 1000, []string{})
	q.RawQuery = `
			SELECT true
			FROM topologies_org_cloud_ctx_ns
			WHERE EXISTS (  SELECT 1 
							FROM topologies_org_cloud_ctx_ns
							WHERE org_id = $1 AND cloud_provider = $2 AND context = $3 AND region = $4 AND namespace = $5
			    		 )`
	return q
}

func (t *InfraBaseTopology) IsOrgCloudCtxNsAuthorized(ctx context.Context, kns kns.TopologyKubeCtxNs) (bool, error) {
	q := getIsOrgCloudCtxNsAuthorizedQueryParams()
	log.Debug().Interface("IsOrgCloudCtxNsAuthorized", q.LogHeader(Sn))
	authorized := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, t.OrgID, kns.CloudProvider, kns.Context, kns.Region, kns.Namespace).Scan(&authorized)
	if err != nil {
		return false, errors.New("not authorized")
	}
	return authorized, err
}

func IsOrgCloudCtxNsAuthorized(ctx context.Context, orgID int, kns kns.TopologyKubeCtxNs) (bool, error) {
	q := getIsOrgCloudCtxNsAuthorizedQueryParams()
	log.Debug().Interface("IsOrgCloudCtxNsAuthorized", q.LogHeader(Sn))
	authorized := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, kns.CloudProvider, kns.Context, kns.Region, kns.Namespace).Scan(&authorized)
	if err != nil {
		return false, errors.New("not authorized")
	}
	return authorized, err
}
