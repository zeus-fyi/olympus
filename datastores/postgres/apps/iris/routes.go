package iris_models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertOrgRoute(ctx context.Context, route iris_autogen_bases.OrgRoutes) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_routes(route_id, org_id, route_path)
				  VALUES ($1, $2, $3)`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, route.RouteID, route.OrgID, route.RoutePath)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRoute"))
}

func InsertOrgRouteGroup(ctx context.Context, ogr iris_autogen_bases.OrgRouteGroups) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_route_groups(route_group_id, org_id, route_group_name)
				  VALUES ($1, $2, $3)`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, ogr.RouteGroupID, ogr.OrgID, ogr.RouteGroupName)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRouteGroup"))
}

func InsertOrgRoutesGroups(ctx context.Context, ors iris_autogen_bases.OrgRoutesGroups) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_routes_groups(route_group_id,route_id)
				  VALUES ($1, $2, $3)`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, ors.RouteGroupID, ors.RouteID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRoutesGroups"))
}

func SelectOrgRoutes(ctx context.Context, orgID int) (iris_autogen_bases.OrgRoutesSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT route_id, org_id, route_path
				  FROM org_routes
				  WHERE org_id = $1`

	var routes iris_autogen_bases.OrgRoutesSlice
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var route iris_autogen_bases.OrgRoutes
		rowErr := rows.Scan(
			&route.RouteID, &route.OrgID, &route.RoutePath,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectOrgRoutes"))
			return nil, rowErr
		}
		routes = append(routes, route)
	}
	return routes, misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes"))
}
