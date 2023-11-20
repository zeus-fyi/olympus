package hera_search

import (
	"context"
	"strconv"

	twitter2 "github.com/cvcio/twitter"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
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
	q.RawQuery = `
        SELECT sq.search_id, sq.query, sq.max_results, COALESCE(MAX(it.tweet_id), 0) AS max_tweet_id
        FROM public.ai_twitter_search_query sq
        LEFT JOIN public.ai_incoming_tweets it ON sq.search_id = it.search_id
        WHERE sq.org_id = $1 AND sq.user_id = $2 AND sq.search_group_name = $3
        GROUP BY sq.search_id, sq.query, sq.max_results;
    `
	return q
}

type TwitterSearchQuery struct {
	SearchID   int    `json:"search_id"`
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
	MaxTweetID int    `json:"max_tweet_id"`
}

func SelectTwitterSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) (*TwitterSearchQuery, error) {
	queryTemplate := selectTwitterSearchQuery()
	ts := &TwitterSearchQuery{}
	err := apps.Pg.QueryRowWArgs(ctx, queryTemplate.RawQuery, ou.OrgID, ou.UserID, searchGroupName).Scan(&ts.SearchID, &ts.Query, &ts.MaxResults, &ts.MaxTweetID)
	if err != nil {
		log.Err(err).Msg("SelectTwitterSearchQuery")
		return nil, err
	}
	return ts, err
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
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tweetIDs, nil
}
