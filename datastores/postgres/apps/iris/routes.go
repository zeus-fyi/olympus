package iris_models

import (
	"context"
	"errors"

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

func InsertOrgRoutesFromQuickNodeID(ctx context.Context, quickNodeID string, routes []iris_autogen_bases.OrgRoutes) error {
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
        WITH cte_qn_org_id AS (
			SELECT ou.org_id as org_id
			FROM quicknode_marketplace_customer qmc 
			LEFT JOIN users_keys usk ON usk.public_key = qmc.quicknode_id
			LEFT JOIN org_users ou ON ou.user_id = usk.user_id
			WHERE quicknode_id = $1
			GROUP BY ou.org_id
			LIMIT 1
		), new_routes (route_id, route_path) AS (
            SELECT * FROM UNNEST ($2::int8[], $3::text[])
        ), existing_routes AS (
            SELECT route_path FROM org_routes WHERE org_id = (SELECT org_id FROM cte_qn_org_id)
        )
        INSERT INTO org_routes (route_id, org_id, route_path)
		SELECT nr.route_id, (SELECT org_id FROM cte_qn_org_id), nr.route_path
		FROM new_routes nr
		WHERE NOT EXISTS (SELECT 1 FROM existing_routes er WHERE er.route_path = nr.route_path)
		ON CONFLICT (org_id, route_path) DO NOTHING
    `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, quickNodeID, pq.Array(routeIDs), pq.Array(routePaths))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No new routes to insert")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRoutesFromQuickNodeID"))
}

func UpsertGeneratedQuickNodeOrgRouteGroup(ctx context.Context, quickNodeID string, ogr iris_autogen_bases.OrgRouteGroups, routes []iris_autogen_bases.OrgRoutes) error {
	// Convert the routes slice into a format that can be used in the SQL query
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.RoutePath
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
      WITH cte_qn_org_id AS (
			SELECT ou.org_id as org_id
			FROM quicknode_marketplace_customer qmc 
			LEFT JOIN users_keys usk ON usk.public_key = qmc.quicknode_id
			LEFT JOIN org_users ou ON ou.user_id = usk.user_id
			WHERE quicknode_id = $1
			GROUP BY ou.org_id
			LIMIT 1
		), cte_upsert_route_group AS (
			INSERT INTO org_route_groups(route_group_id, org_id, route_group_name, auto_generated)
			VALUES ($2, (SELECT org_id FROM cte_qn_org_id), $3, true)
			ON CONFLICT (org_id, route_group_name) DO UPDATE SET 
				auto_generated = EXCLUDED.auto_generated
			RETURNING route_group_id
		), cte_route_ids AS (
			SELECT route_id as route_id
			FROM org_routes
			WHERE org_id = (SELECT org_id FROM cte_qn_org_id) AND route_path = ANY($4::text[])
		) 	  INSERT INTO org_routes_groups(route_id, route_group_id)
			  SELECT route_id, (SELECT COALESCE(route_group_id, $2) FROM cte_upsert_route_group) as route_group_id
			  FROM cte_route_ids
			  ON CONFLICT (route_id, route_group_id) DO NOTHING
	`
	ogr.RouteGroupID = ts.UnixTimeStampNow()
	_, err := apps.Pg.Exec(ctx, q.RawQuery, quickNodeID, ogr.RouteGroupID, ogr.RouteGroupName, pq.Array(routePaths))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No new routes to insert")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgRouteGroup"))
}

func InsertOrgRouteGroup(ctx context.Context, ogr iris_autogen_bases.OrgRouteGroups, routes []iris_autogen_bases.OrgRoutes) error {
	// Convert the routes slice into a format that can be used in the SQL query
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.RoutePath
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		WITH cte_entry AS (
			SELECT route_id
	 		FROM org_routes
			WHERE org_id = $2 AND route_path = ANY($4::text[])
 		), cte_del AS (
			DELETE FROM org_routes_groups
			WHERE route_group_id IN (SELECT route_group_id FROM org_route_groups WHERE org_id = $2 AND route_group_name = $3)
		), new_route_group AS (
			INSERT INTO org_route_groups(route_group_id, org_id, route_group_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (org_id, route_group_name) DO UPDATE SET route_group_id = $1
			RETURNING route_group_id
		), cte_ins AS (
			SELECT COALESCE(route_group_id, $1) as route_group_id
			FROM new_route_group
		) INSERT INTO org_routes_groups(route_id, route_group_id)
		  SELECT route_id, (SELECT COALESCE(route_group_id, $1) FROM cte_ins) as route_group_id
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
	Map map[int]map[string][]RouteInfo
}

func SelectAllOrgRoutes(ctx context.Context) (OrgRoutesGroup, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]RouteInfo),
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT 
					COALESCE(orgg.route_group_name, 'unused') as route_group_name,
					org.route_path,
					org.org_id,
					NULLIF(array_remove(array_agg(r.referer), NULL)::text[], ARRAY[]::text[]) as referers
				FROM 
					org_routes org
				LEFT JOIN
					org_routes_groups orgrs ON org.route_id = orgrs.route_id
				LEFT JOIN
					org_route_groups orgg ON orgg.route_group_id = orgrs.route_group_id
				LEFT JOIN 
					provisioned_quicknode_services pqs ON org.route_path = pqs.http_url
				LEFT JOIN 
					provisioned_quicknode_services_referers r ON pqs.endpoint_id = r.endpoint_id
				GROUP BY 
					org.org_id,
					orgg.route_group_name,
					org.route_path;
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutes")); returnErr != nil {
		return og, err
	}
	defer rows.Close()
	for rows.Next() {
		var orgID int
		var routeGroupName string
		var routePath string
		var referers []string

		rowErr := rows.Scan(&routeGroupName, &routePath, &orgID, pq.Array(&referers))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectAllOrgRoutes"))
			return og, rowErr
		}

		if _, ok := og.Map[orgID]; !ok {
			og.Map[orgID] = make(map[string][]RouteInfo)
		}

		og.Map[orgID][routeGroupName] = append(og.Map[orgID][routeGroupName], RouteInfo{
			RoutePath: routePath,
			Referers:  referers,
		})
	}
	return og, misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutes"))
}

type RouteInfo struct {
	RoutePath string
	Referers  []string
}

func SelectAllOrgRoutesByOrg(ctx context.Context, orgID int) (map[string][]RouteInfo, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]RouteInfo),
	}

	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT 
					COALESCE(orgg.route_group_name, 'unused') as route_group_name,
					org.route_path,
					NULLIF(array_remove(array_agg(r.referer), NULL)::text[], ARRAY[]::text[]) as referers
					FROM 
						org_routes org
					LEFT JOIN
						org_routes_groups orgrs ON org.route_id = orgrs.route_id
					LEFT JOIN
						org_route_groups orgg ON orgg.route_group_id = orgrs.route_group_id
					LEFT JOIN 
						provisioned_quicknode_services pqs ON org.route_path = pqs.http_url
					LEFT JOIN 
						provisioned_quicknode_services_referers r ON pqs.endpoint_id = r.endpoint_id
					WHERE org.org_id = $1
					GROUP BY 
						orgg.route_group_name,
						org.route_path;
      `

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutesByOrg")); returnErr != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var routeGroupName string
		var routePath string
		var referers []string

		rowErr := rows.Scan(&routeGroupName, &routePath, pq.Array(&referers))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectAllOrgRoutesByOrg"))
			return nil, rowErr
		}

		if _, ok := og.Map[orgID]; !ok {
			og.Map[orgID] = make(map[string][]RouteInfo)
		}

		og.Map[orgID][routeGroupName] = append(og.Map[orgID][routeGroupName], RouteInfo{
			RoutePath: routePath,
			Referers:  referers,
		})
	}

	return og.Map[orgID], misc.ReturnIfErr(err, q.LogHeader("SelectAllOrgRoutesByOrg"))
}

type OrgRoutesGroupsAndEndpoints struct {
	Map    map[string][]string `json:"orgGroupRoutes"`
	Routes []string            `json:"routes"`
}

func SelectAllEndpointsAndOrgGroupRoutesByOrg(ctx context.Context, orgID int) (OrgRoutesGroupsAndEndpoints, error) {
	og := OrgRoutesGroupsAndEndpoints{
		Map:    make(map[string][]string),
		Routes: []string{},
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE(org.route_group_name, 'unused'), o.org_id, o.route_path
				  FROM org_routes o 
				  LEFT JOIN org_routes_groups orgrs ON orgrs.route_id = o.route_id
				  LEFT JOIN org_route_groups org ON org.route_group_id = orgrs.route_group_id
				  WHERE o.org_id = $1
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectAllEndpointsAndOrgGroupRoutesByOrg")); returnErr != nil {
		return og, err
	}
	seenMap := make(map[string]bool)

	defer rows.Close()
	for rows.Next() {
		var route iris_autogen_bases.OrgRoutes
		gn := ""
		rowErr := rows.Scan(
			&gn, &route.OrgID, &route.RoutePath,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectAllEndpointsAndOrgGroupRoutesByOrg"))
			return og, rowErr
		}
		if _, ok := og.Map[gn]; !ok {
			og.Map[gn] = []string{route.RoutePath}
		} else {
			tmp := og.Map[gn]
			tmp = append(tmp, route.RoutePath)
			og.Map[gn] = tmp
		}
		if seenMap[route.RoutePath] != true {
			og.Routes = append(og.Routes, route.RoutePath)
			seenMap[route.RoutePath] = true
		}
	}
	return og, misc.ReturnIfErr(err, q.LogHeader("SelectAllEndpointsAndOrgGroupRoutesByOrg"))
}

func SelectOrgRoutesByOrgAndGroupName(ctx context.Context, orgID int, groupName string) (OrgRoutesGroup, error) {
	og := OrgRoutesGroup{
		Map: make(map[int]map[string][]RouteInfo),
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT 
					o.route_group_name, 
					o.org_id, 
					org.route_path, 
					NULLIF(array_remove(array_agg(r.referer), NULL)::text[], ARRAY[]::text[]) as referers
				  FROM org_route_groups o 
				  INNER JOIN org_routes_groups orgrs ON orgrs.route_group_id = o.route_group_id
				  LEFT JOIN org_routes org ON org.route_id = orgrs.route_id
				  LEFT JOIN 
					provisioned_quicknode_services pqs ON org.route_path = pqs.http_url
				  LEFT JOIN 
					provisioned_quicknode_services_referers r ON pqs.endpoint_id = r.endpoint_id
				  WHERE o.org_id = $1 AND o.route_group_name = $2
				  GROUP BY 
						o.route_group_name,
						o.org_id,
						org.route_path;
				  `

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, groupName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectOrgRoutes")); returnErr != nil {
		return og, err
	}
	defer rows.Close()
	for rows.Next() {
		var routeGroupName string
		var routePath string
		var referers []string
		rowErr := rows.Scan(
			&routeGroupName, &orgID, &routePath, pq.Array(&referers),
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectOrgRoutes"))
			return og, rowErr
		}
		if _, ok := og.Map[orgID]; !ok {
			og.Map[orgID] = make(map[string][]RouteInfo)
		}
		og.Map[orgID][routeGroupName] = append(og.Map[orgID][routeGroupName], RouteInfo{
			RoutePath: routePath,
			Referers:  referers,
		})
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

func DeleteOrgGroupAndRoutes(ctx context.Context, orgID int, routeGroupName string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		WITH cte_delete1 AS (
			DELETE FROM org_routes_groups og
		    WHERE route_group_id IN ( SELECT ortg.route_group_id
			  						  FROM org_routes_groups ortg
			  						  INNER JOIN org_route_groups org ON org.route_group_id = ortg.route_group_id		
									  WHERE org.org_id = $1 AND org.route_group_name = $2)
		) 
		DELETE FROM org_route_groups
		WHERE org_id = $1 AND route_group_name = $2
	`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, routeGroupName)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No routes to delete")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeleteOrgGroupAndRoutes"))
}

type TableUsage struct {
	EndpointCount int `json:"endpointCount"`
	TableCount    int `json:"tableCount"`
}

func OrgEndpointsAndGroupTablesCount(ctx context.Context, orgID int) (*TableUsage, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		SELECT COALESCE(COUNT(*), 0) as endpoint_count, 
       		COALESCE(
       		   (SELECT COUNT(*)
       		    FROM org_route_groups WHERE org_id = $1  AND auto_generated = false
       		    AND EXISTS (SELECT 1 FROM org_routes_groups WHERE org_routes_groups.route_group_id = org_route_groups.route_group_id))
       		    ,0) as table_count
		FROM org_routes 
		WHERE org_id = $1
	`

	endpointCount, groupTablesCount := 0, 0
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID).Scan(&endpointCount, &groupTablesCount)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("OrgEndpointsAndGroupTablesCount has no entries")
		return &TableUsage{0, 0}, nil
	}
	if err != nil {
		return nil, err
	}

	return &TableUsage{endpointCount, groupTablesCount}, misc.ReturnIfErr(err, q.LogHeader("OrgEndpointsAndGroupTablesCount"))
}

func OrgGroupTablesToRemove(ctx context.Context, quickNodeID string, plan string) ([]string, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		SELECT route_group_id, route_group_name
		FROM org_route_groups
		WHERE org_id = (SELECT org_id FROM users_keys usk
						JOIN org_users ou ON ou.user_id = usk.user_id
						WHERE public_key = $1
						LIMIT 1)
		AND EXISTS (SELECT 1 FROM org_routes_groups WHERE org_routes_groups.route_group_id = org_route_groups.route_group_id)
		AND auto_generated = false
		ORDER BY route_group_id
	`

	maxCount := 1
	switch plan {
	case "performance":
		maxCount = PerformanceGroupTables
	case "standard":
		maxCount = StandardGroupTables
	case "test":
		maxCount = FreeGroupTables
	case "free":
		maxCount = FreeGroupTables
	}

	rows, err := apps.Pg.Query(ctx, q.RawQuery, quickNodeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("OrgGroupTablesToRemove")); returnErr != nil {
		return nil, err
	}

	var ogToDelete []string
	count := 0
	defer rows.Close()
	for rows.Next() {
		var routeGroupName string
		var routeGroupID int

		rowErr := rows.Scan(
			&routeGroupID, &routeGroupName,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("OrgGroupTablesToRemove"))
			return nil, rowErr
		}
		if count >= maxCount {
			ogToDelete = append(ogToDelete, routeGroupName)
		}
		count += 1
	}
	return ogToDelete, misc.ReturnIfErr(err, q.LogHeader("OrgGroupTablesToRemove"))
}

/*
SELECT route_group_id, route_group_name
FROM org_route_groups
WHERE org_id = 7138983863666903883
AND EXISTS (SELECT 1 FROM org_routes_groups WHERE org_routes_groups.route_group_id = org_route_groups.route_group_id)
ORDER BY route_group_id
*/
const (
	FreeGroupTables        = 1
	StandardGroupTables    = 50
	PerformanceGroupTables = 250
)

func (t *TableUsage) CheckEndpointLimits() error {
	if t.EndpointCount > 1000 {
		return errors.New("exceeds plan endpoints")
	}
	return nil
}

func (t *TableUsage) CheckPlanLimits(plan string) error {
	err := t.CheckEndpointLimits()
	if err != nil {
		log.Err(err).Msg("CheckPlanLimits")
		return err
	}
	switch plan {
	case "performance":
		// check 50k ZU/s
		// check max 3B ZU/month
		if t.TableCount >= PerformanceGroupTables {
			return errors.New("exceeds plan group tables")
		}
		return nil
	case "standard":
		// check 25k ZU/s
		// check max 1B ZU/month
		if t.TableCount >= StandardGroupTables {
			return errors.New("exceeds plan group tables")
		}
		return nil
	case "free":
		// check 1k ZU/s
		// check max 50M ZU/month
		if t.TableCount >= FreeGroupTables {
			return errors.New("exceeds plan group tables")
		}
		return nil
	case "test":
	default:
	}
	return nil
}
