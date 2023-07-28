package iris_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
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

var ts chronos.Chronos

func InsertOrgRoutes(ctx context.Context, orgID int, routes []iris_autogen_bases.OrgRoutes) error {
	// Generate a slice of IDs for the new routes
	routeIDs := make([]int, len(routes))
	for i := range routeIDs {
		routeIDs[i] = ts.UnixTimeStampNow()
	}

	// Convert the routes slice into a format that can be used in the SQL query
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.RoutePath
	}

	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        WITH new_routes (route_id, route_path) AS (
            SELECT * FROM UNNEST ($2::int8[], $3::text[])
        ), existing_routes AS (
            SELECT route_path FROM org_routes WHERE org_id = $1
        )
        INSERT INTO org_routes (route_id, org_id, route_path)
		SELECT nr.route_id, $1, nr.route_path
		FROM new_routes nr
		WHERE NOT EXISTS (SELECT 1 FROM existing_routes er WHERE er.route_path = nr.route_path)
    `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, pq.Array(routeIDs), pq.Array(routePaths))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No new routes to insert")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRoutes"))
}

func InsertOrgRouteGroup(ctx context.Context, ogr iris_autogen_bases.OrgRouteGroups, routes []iris_autogen_bases.OrgRoutes) error {
	// Convert the routes slice into a format that can be used in the SQL query
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.RoutePath
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		WITH new_route_group AS (
			INSERT INTO org_route_groups(route_group_id, org_id, route_group_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (org_id, route_group_name) DO NOTHING
			RETURNING route_group_id
		), cte_entry AS (
			SELECT route_id
	 		FROM org_routes
			WHERE org_id = $2 AND route_path = ANY($4::text[])
 		), cte_rg_id AS (
			SELECT route_group_id 
			FROM org_route_groups
			WHERE org_id = $2 AND route_group_name = $3
		), cte_del AS (
			DELETE FROM org_routes_groups
			WHERE route_id NOT IN (SELECT route_id FROM cte_entry)
		) INSERT INTO org_routes_groups(route_id, route_group_id)
		  SELECT route_id, (SELECT route_group_id FROM cte_rg_id) as route_group_id
		  FROM cte_entry
		  ON CONFLICT (route_id, route_group_id) DO NOTHING
	`
	ogr.RouteGroupID = ts.UnixTimeStampNow()
	_, err := apps.Pg.Exec(ctx, q.RawQuery, ogr.RouteGroupID, ogr.OrgID, ogr.RouteGroupName, pq.Array(routePaths))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No new routes to insert")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRouteGroup"))
}

func InsertOrgRoutesGroups(ctx context.Context, ors iris_autogen_bases.OrgRoutesGroups) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_routes_groups(route_group_id,route_id)
				  VALUES ($1, $2)`
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

type OrgRoutesGroup struct {
	Map map[int]map[string][]string
}

func SelectAllOrgRoutes(ctx context.Context) (OrgRoutesGroup, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]string),
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT o.route_group_name, o.org_id, org.route_path
				  FROM org_route_groups o 
				  INNER JOIN org_routes_groups orgrs ON orgrs.route_group_id = o.route_group_id
				  LEFT JOIN org_routes org ON org.route_id = orgrs.route_id
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutes")); returnErr != nil {
		return og, err
	}
	defer rows.Close()
	for rows.Next() {
		var route iris_autogen_bases.OrgRoutes
		gn := ""
		rowErr := rows.Scan(
			&gn, &route.OrgID, &route.RoutePath,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectAllOrgRoutes"))
			return og, rowErr
		}
		if _, ok := og.Map[route.OrgID]; !ok {
			og.Map[route.OrgID] = make(map[string][]string)
		}
		if _, ok := og.Map[route.OrgID][gn]; !ok {
			og.Map[route.OrgID][gn] = []string{}
		}
		tmp := og.Map[route.OrgID][gn]
		tmp = append(tmp, route.RoutePath)
		og.Map[route.OrgID][gn] = tmp
	}
	return og, misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutes"))
}

func SelectAllOrgRoutesByOrg(ctx context.Context, orgID int) (OrgRoutesGroup, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]string),
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT o.route_group_name, o.org_id, org.route_path
				  FROM org_route_groups o 
				  INNER JOIN org_routes_groups orgrs ON orgrs.route_group_id = o.route_group_id
				  LEFT JOIN org_routes org ON org.route_id = orgrs.route_id
				  WHERE o.org_id = $1
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes")); returnErr != nil {
		return og, err
	}
	defer rows.Close()
	for rows.Next() {
		var route iris_autogen_bases.OrgRoutes
		gn := ""
		rowErr := rows.Scan(
			&gn, &route.OrgID, &route.RoutePath,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectOrgRoutes"))
			return og, rowErr
		}
		if _, ok := og.Map[route.OrgID]; !ok {
			og.Map[route.OrgID] = make(map[string][]string)
		}
		if _, ok := og.Map[route.OrgID][gn]; !ok {
			og.Map[route.OrgID][gn] = []string{}
		}
		tmp := og.Map[route.OrgID][gn]
		tmp = append(tmp, route.RoutePath)
		og.Map[route.OrgID][gn] = tmp
	}
	return og, misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes"))
}

func SelectOrgRoutesByOrgAndGroupName(ctx context.Context, orgID int, groupName string) (OrgRoutesGroup, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]string),
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT o.route_group_name, o.org_id, org.route_path
				  FROM org_route_groups o 
				  INNER JOIN org_routes_groups orgrs ON orgrs.route_group_id = o.route_group_id
				  LEFT JOIN org_routes org ON org.route_id = orgrs.route_id
				  WHERE o.org_id = $1 AND o.route_group_name = $2
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, groupName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes")); returnErr != nil {
		return og, err
	}
	defer rows.Close()
	for rows.Next() {
		var route iris_autogen_bases.OrgRoutes
		gn := ""
		rowErr := rows.Scan(
			&gn, &route.OrgID, &route.RoutePath,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectOrgRoutes"))
			return og, rowErr
		}
		if _, ok := og.Map[route.OrgID]; !ok {
			og.Map[route.OrgID] = make(map[string][]string)
		}
		if _, ok := og.Map[route.OrgID][gn]; !ok {
			og.Map[route.OrgID][gn] = []string{}
		}
		tmp := og.Map[route.OrgID][gn]
		tmp = append(tmp, route.RoutePath)
		og.Map[route.OrgID][gn] = tmp
	}
	return og, misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes"))
}

func DeleteOrgRoutes(ctx context.Context, orgID int, routes []string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		WITH cte_delete1 AS (
			DELETE FROM org_routes_groups og
		    WHERE route_id IN (SELECT route_id FROM org_routes WHERE org_id = $1 AND route_path = ANY($2::text[]))
		) DELETE FROM org_routes
		  WHERE route_id IN (SELECT route_id FROM org_routes WHERE org_id = $1 AND route_path = ANY($2::text[]))
	`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, pq.Array(routes))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No routes to delete")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeleteOrgRoutes"))
}
