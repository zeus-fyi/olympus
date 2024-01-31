package iris_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func DeleteOrgRoutesFromGroup(ctx context.Context, orgID int, groupName string, routePaths []string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `DELETE FROM org_routes_groups ogs
				  USING org_routes orr, org_route_groups orgg
				  WHERE orr.route_id = ogs.route_id
				  AND orgg.route_group_id = ogs.route_group_id
				  AND orr.org_id = $1
				  AND orgg.route_group_name = $2
				  AND orr.route_path IN (SELECT UNNEST($3::bigint[]))`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, groupName, pq.Array(routePaths))
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No new routes to insert")
		return nil
	}
	if err != nil {
		log.Err(err).Int("orgID", orgID).Str("groupName", groupName).Interface("routes", routePaths).Msg("DeleteOrgRoutesFromGroup")
		return err
	}
	return err
}

func DeleteOrgRoutingGroup(ctx context.Context, orgID int, groupName string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
		WITH cte_entry AS (
			  DELETE FROM org_routes_groups ogs
			  USING org_routes orr, org_route_groups orgg
			  WHERE orr.route_id = ogs.route_id
			  AND orgg.route_group_id = ogs.route_group_id
			  AND orr.org_id = $1
			  AND orgg.route_group_name = $2
		) DELETE FROM org_route_groups 
 		  WHERE route_group_name = $2 AND org_id = $1`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, groupName)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("no routes to delete")
		return nil
	}
	if err != nil {
		log.Err(err).Int("orgID", orgID).Str("groupName", groupName).Msg("DeleteOrgRoutingGroup")
		return err
	}
	return err
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
	if err != nil {
		log.Err(err).Int("orgID", orgID).Interface("routes", routes).Msg("DeleteOrgRoutes")
		return err
	}
	return err
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
	case "enterprise":
		maxCount = EnterpriseGroupTables
	case "performance":
		maxCount = PerformanceGroupTables
	case "standard":
		maxCount = StandardGroupTables
	case "lite":
		maxCount = LiteGroupTables
	case "test":
		maxCount = FreeGroupTables
	case "free":
		maxCount = FreeGroupTables
	}

	rows, err := apps.Pg.Query(ctx, q.RawQuery, quickNodeID)
	if err != nil {
		log.Err(err).Str("quickNodeID", quickNodeID).Interface("plan", plan).Msg("OrgGroupTablesToRemove")
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
			log.Err(rowErr).Str("quickNodeID", quickNodeID).Interface("plan", plan).Msg("OrgGroupTablesToRemove")
			return nil, rowErr
		}
		if count >= maxCount {
			ogToDelete = append(ogToDelete, routeGroupName)
		}
		count += 1
	}
	if err != nil {
		log.Err(err).Str("quickNodeID", quickNodeID).Interface("plan", plan).Msg("OrgGroupTablesToRemove")
		return nil, err
	}
	return ogToDelete, err
}

//func DeleteOrgGroupAndRoutes(ctx context.Context, orgID int, routeGroupName string) error {
//	q := sql_query_templates.QueryParams{}
//	q.RawQuery = `
//		WITH cte_delete1 AS (
//			DELETE FROM org_routes_groups og
//		    WHERE route_group_id IN ( SELECT ortg.route_group_id
//			  						  FROM org_routes_groups ortg
//			  						  INNER JOIN org_route_groups org ON org.route_group_id = ortg.route_group_id
//									  WHERE org.org_id = $1 AND org.route_group_name = $2)
//		)
//		DELETE FROM org_route_groups
//		WHERE org_id = $1 AND route_group_name = $2
//	`
//	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, routeGroupName)
//	if err == pgx.ErrNoRows {
//		log.Warn().Msg("No routes to delete")
//		return nil
//	}
//	return misc.ReturnIfErr(err, q.LogHeader("DeleteOrgGroupAndRoutes"))
//}
