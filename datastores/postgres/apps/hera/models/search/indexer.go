package hera_search

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func UpdateDiscordSearchQueryStatus(ctx context.Context, ou org_users.OrgUser, sp SearchIndexerParams) error {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "updateDiscordSearchQueryStatus"
	q.RawQuery = `
        UPDATE "public"."ai_discord_search_query"
        SET "active" = $1
        WHERE "org_id" = $2 AND "search_group_name" = $3;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, sp.Active, ou.OrgID, sp.SearchGroupName)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("UpdateDiscordSearchQueryStatus")
		return err
	}
	return nil
}

func UpdateRedditSearchQueryStatus(ctx context.Context, ou org_users.OrgUser, sp SearchIndexerParams) error {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "UpdateRedditSearchQueryStatus"
	q.RawQuery = `
        UPDATE "public"."ai_reddit_search_query"
        SET "active" = $1
        WHERE "org_id" = $2 AND "search_group_name" = $3;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, sp.Active, ou.OrgID, sp.SearchGroupName)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("UpdateRedditSearchQueryStatus")
		return err
	}
	return nil
}

func UpdateTwitterSearchQueryStatus(ctx context.Context, ou org_users.OrgUser, sp SearchIndexerParams) error {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "UpdateTwitterSearchQueryStatus"
	q.RawQuery = `
        UPDATE "public"."ai_twitter_search_query"
        SET "active" = $1
        WHERE "org_id" = $2 AND "search_group_name" = $3;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, sp.Active, ou.OrgID, sp.SearchGroupName)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("UpdateTwitterSearchQueryStatus")
		return err
	}
	return nil
}

//func UpdateTelegramSearchQueryStatus(ctx context.Context, ou org_users.OrgUser, searchGroupName string, active bool) error {
//	q := sql_query_templates.QueryParams{}
//	q.QueryName = "UpdateTelegramSearchQueryStatus"
//	q.RawQuery = `
//        UPDATE "public".""
//        SET "active" = $1
//        WHERE "org_id" = $2 AND "search_group_name" = $3;`
//
//	_, err := apps.Pg.Exec(ctx, q.RawQuery, active, ou.OrgID, searchGroupName)
//	if err != nil && err != pgx.ErrNoRows {
//		log.Err(err).Msg("UpdateTelegramSearchQueryStatus")
//		return err
//	}
//	return nil
//}
