package read_topology

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/infra"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology `json:"infraBaseTopology"`
	read_charts.Chart       `json:"chart"`
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

func (t *InfraBaseTopology) SelectInfraTopologyQueryForOrg() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQueryForOrg"
	q.CTEQuery.Params = append(q.CTEQuery.Params, t.TopologyID, t.OrgID)
	q.RawQuery = read_charts.FetchChartQuery(q)
	return q
}

func (t *InfraBaseTopology) SelectTopologyForOrg(ctx context.Context) error {
	q := t.SelectInfraTopologyQueryForOrg()

	log.Debug().Interface("SelectTopologyQuery", q.LogHeader(Sn))
	err := t.SelectSingleChartsResources(ctx, q)
	if err != nil {
		log.Err(err).Interface("topology", t).Msg("SelectTopology, SelectSingleChartsResources error")
		return err
	}
	return err
}

func (t *InfraBaseTopology) SelectTopology(ctx context.Context) error {
	q := t.SelectInfraTopologyQueryForOrg()

	log.Debug().Interface("SelectInfraTopologyQueryForOrg", q.LogHeader(Sn))
	err := t.SelectSingleChartsResources(ctx, q)
	if err != nil {
		log.Err(err).Interface("topology", t).Msg("SelectTopology, SelectSingleChartsResources error")
		return err
	}
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

func getIsOrgCloudCtxNsAuthorizedQueryParamsFromID() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("IsOrgCloudCtxNsAuthorized", "topologies_org_cloud_ctx_ns", "where", 1000, []string{})
	q.RawQuery = `
			SELECT true, cloud_provider, region, context, namespace
			FROM topologies_org_cloud_ctx_ns
			WHERE EXISTS (  SELECT 1 
							FROM topologies_org_cloud_ctx_ns
							WHERE org_id = $1 AND cloud_ctx_ns_id = $2 
			    		 ) AND cloud_ctx_ns_id = $2`
	return q
}
func IsOrgCloudCtxNsAuthorizedFromID(ctx context.Context, orgID, cloudCtxNsID int) (bool, zeus_common_types.CloudCtxNs, error) {
	q := getIsOrgCloudCtxNsAuthorizedQueryParamsFromID()
	log.Debug().Interface("IsOrgCloudCtxNsAuthorizedFromID", q.LogHeader(Sn))
	authorized := false
	cctx := zeus_common_types.CloudCtxNs{}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, cloudCtxNsID).Scan(&authorized, &cctx.CloudProvider, &cctx.Region, &cctx.Context, &cctx.Namespace)
	if err != nil {
		if orgID == TemporalOrgID {
			log.Ctx(ctx).Info().Msg("IsOrgCloudCtxNsAuthorized: Using Temporal Key")
			return true, cctx, nil
		}
		return false, cctx, errors.New("not authorized")
	}
	return authorized, cctx, err
}

const TemporalOrgID = 7138983863666903883

func (t *InfraBaseTopology) IsOrgCloudCtxNsAuthorized(ctx context.Context, kns zeus_common_types.CloudCtxNs) (bool, error) {
	q := getIsOrgCloudCtxNsAuthorizedQueryParams()
	log.Debug().Interface("IsOrgCloudCtxNsAuthorized", q.LogHeader(Sn))
	authorized := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, t.OrgID, kns.CloudProvider, kns.Context, kns.Region, kns.Namespace).Scan(&authorized)
	if err != nil {
		if t.OrgID == TemporalOrgID {
			log.Info().Msg("IsOrgCloudCtxNsAuthorized: Using Temporal Key")
			return true, nil
		}
		return false, errors.New("not authorized")
	}
	return authorized, err
}
func IsOrgCloudCtxNsAuthorized(ctx context.Context, orgID int, kns zeus_common_types.CloudCtxNs) (bool, error) {
	q := getIsOrgCloudCtxNsAuthorizedQueryParams()
	log.Debug().Interface("IsOrgCloudCtxNsAuthorized", q.LogHeader(Sn))
	authorized := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, kns.CloudProvider, kns.Context, kns.Region, kns.Namespace).Scan(&authorized)
	if err != nil {
		if orgID == TemporalOrgID {
			log.Info().Msg("IsOrgCloudCtxNsAuthorized: Using Temporal Key")
			return true, nil
		}
		return false, errors.New("not authorized")
	}
	return authorized, err
}
