package hera_search

import (
	"context"
	"strconv"

	twitter2 "github.com/cvcio/twitter"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const (
	defaultTwitterSearchGroupName = "zeusfyi"
)

func insertTwitterSearchQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertTwitterSearchQuery"
	q.RawQuery = `INSERT INTO "public"."ai_twitter_search_query" ("org_id", "user_id", "search_group_name", "max_results", "query")
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "search_id";`
	return q
}

func InsertTwitterSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string, query string, maxResults int) (int, error) {
	queryTemplate := insertTwitterSearchQuery()
	var searchID int
	err := apps.Pg.QueryRowWArgs(ctx, queryTemplate.RawQuery, ou.OrgID, ou.UserID, searchGroupName, maxResults, query).Scan(&searchID)
	if err != nil {
		return 0, err
	}
	return searchID, nil
}

func selectTwitterSearchQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectTwitterSearchQuery"
	q.RawQuery = `SELECT search_id, query, max_results
				  FROM ai_twitter_search_query 
				  WHERE org_id = $1 AND user_id = $2 AND search_group_name = $3;`
	return q
}

func SelectTwitterSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) (int, string, int, error) {
	queryTemplate := selectTwitterSearchQuery()
	var searchID int
	var maxResults int
	var query string
	err := apps.Pg.QueryRowWArgs(ctx, queryTemplate.RawQuery, ou.OrgID, ou.UserID, searchGroupName).Scan(&searchID, &query, &maxResults)
	if err != nil {
		return 0, query, maxResults, err
	}
	return searchID, query, maxResults, err
}

func insertIncomingTweets() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertIncomingTweets"
	q.RawQuery = `INSERT INTO "public"."ai_incoming_tweets" ("search_id", "tweet_id", "message_text")
        VALUES ($1, $2, $3)
        ON CONFLICT ("tweet_id")
        DO UPDATE SET
            "message_text" = EXCLUDED."message_text"
        RETURNING "tweet_id";`
	return q
}

func InsertIncomingTweets(ctx context.Context, searchID int, tweets []*twitter2.Tweet) ([]int, error) {
	queryTemplate := insertIncomingTweets()
	var tweetIDs []int

	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	for _, tweet := range tweets {
		if tweet == nil {
			continue
		}
		if tweet.Text == "" {
			continue
		}
		var tweetID int
		atoi, aerr := strconv.Atoi(tweet.ID)
		if aerr != nil {
			return nil, aerr
		}
		err = tx.QueryRow(ctx, queryTemplate.RawQuery, searchID, atoi, tweet.Text).Scan(&tweetID)
		if err != nil {
			return nil, err
		}
		if tweetID == 0 {
			continue
		}
		tweetIDs = append(tweetIDs, tweetID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return tweetIDs, nil
}
