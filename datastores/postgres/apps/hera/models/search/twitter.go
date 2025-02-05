package hera_search

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	twitter2 "github.com/cvcio/twitter"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const (
	defaultTwitterSearchGroupName = "zeusfyi"
)

// links are mostly spam, filtering out links reduces the number of tweets by 75%, which results in ~90% less spam

func twitterSearchQuery2() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "twitterSearchQuery"
	q.RawQuery = `SELECT tweet_id, message_text
				  FROM public.ai_incoming_tweets
				  JOIN ai_twitter_search_query sq ON sq.search_id = ai_incoming_tweets.search_id
				  WHERE NOT EXISTS (
					  SELECT 1
					  FROM unnest(ARRAY['🧰','⏳','💥','📍', '🎤', '🚀', '🛑','🏆','🚨','📅','☸️','🆕', '🏓 ', '⭕️','🛡️','👉', '🎟️', '💎', '🪂']) as t(emoji)
					  WHERE message_text LIKE '%' || t.emoji || '%'
						OR (LENGTH(message_text) - LENGTH(REPLACE(message_text, '@', ''))) > 7
						OR (LENGTH(message_text) - LENGTH(REPLACE(message_text, '#', ''))) > 2
					)
				 AND message_text NOT LIKE '%https://t.co%' AND sq.org_id = $1`
	return q
}

func twitterSearchQuery(ou org_users.OrgUser, sp AiSearchParams) (sql_query_templates.QueryParams, []interface{}) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "twitterSearchQuery"
	args := []interface{}{ou.OrgID}

	bq := `SELECT tweet_id, message_text
		   FROM public.ai_incoming_tweets
		   JOIN ai_twitter_search_query sq ON sq.search_id = ai_incoming_tweets.search_id
		   WHERE sq.org_id = $1
		   `
	// Positive keywords
	if sp.Retrieval.RetrievalKeywords != nil && *sp.Retrieval.RetrievalKeywords != "" {
		posQuery := formatKeywordsForTsQuery(*sp.Retrieval.RetrievalKeywords)
		if posQuery != "" {
			bq += fmt.Sprintf(` AND message_text_tsvector @@ to_tsquery('english', $%d)`, len(args)+1)
			args = append(args, posQuery)
		}
	}

	if sp.Retrieval.RetrievalNegativeKeywords != nil && *sp.Retrieval.RetrievalNegativeKeywords != "" {
		negQuery := formatKeywordsForTsQuery(*sp.Retrieval.RetrievalNegativeKeywords, true)
		if negQuery != "" {
			bq += fmt.Sprintf(` AND message_text_tsvector @@ to_tsquery('english', $%d)`, len(args)+1)
			args = append(args, negQuery)
		}
	}
	bq += ` AND NOT EXISTS (
				SELECT 1
				FROM unnest(ARRAY['🧰','⏳','💥','📍', '🎤', '🚀', '🛑','🏆','🚨','📅','☸️','🆕', '🏓 ', '⭕️','🛡️','👉', '🎟️', '💎', '🪂']) as t(emoji)
				WHERE message_text LIKE '%' || t.emoji || '%'
					OR (LENGTH(message_text) - LENGTH(REPLACE(message_text, '@', ''))) > 7
					OR (LENGTH(message_text) - LENGTH(REPLACE(message_text, '#', ''))) > 2
				)`

	if !sp.Window.IsWindowEmpty() {
		bq += ` AND`
		tsRangeStart, tsEnd := sp.Window.GetUnixTimestamps()
		bq += fmt.Sprintf(` tweet_id BETWEEN $%d AND $%d`, len(args)+1, len(args)+2)
		cts := chronos.Chronos{}
		tweetStart := cts.ConvertUnixTimestampToTweetID(tsRangeStart)
		tweetEnd := cts.ConvertUnixTimestampToTweetID(tsEnd)
		args = append(args, tweetStart, tweetEnd)
	}
	bq += ` ORDER BY tweet_id DESC`
	q.RawQuery = bq
	return q, args
}

// Helper function to format keywords for tsquery
func formatKeywordsForTsQuery(queryText string, negate ...bool) string {
	keywords := strings.FieldsFunc(queryText, func(r rune) bool {
		return r == ',' || r == ' '
	})
	formattedQuery := ""
	for _, keyword := range keywords {
		if formattedQuery != "" {
			formattedQuery += " & "
		}
		if len(negate) > 0 && negate[0] {
			formattedQuery += "!" + keyword
		} else {
			formattedQuery += keyword
		}
	}
	return formattedQuery
}

func SearchTwitter(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q, args := twitterSearchQuery(ou, sp)
	var srs []SearchResult
	var rows pgx.Rows
	var err error
	rows, err = apps.Pg.Query(ctx, q.RawQuery, args...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SearchTwitter")); returnErr != nil {
		return nil, err
	}
	ts := chronos.Chronos{}
	defer rows.Close()
	for rows.Next() {
		var sr SearchResult
		sr.Source = "twitter"
		sr.Verified = aws.Bool(true)
		rowErr := rows.Scan(
			&sr.UnixTimestamp, &sr.Value,
		)
		sr.TwitterMetadata = &TwitterMetadata{
			TweetStrID: fmt.Sprintf("%d", sr.UnixTimestamp),
			Text:       sr.Value,
		}
		sr.UnixTimestamp = ts.ConvertTweetIDToUnixTimestamp(sr.UnixTimestamp)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SearchTwitter"))
			return nil, rowErr
		}
		srs = append(srs, sr)
	}
	return srs, nil
}

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
        WHERE sq.org_id = $1 AND sq.search_group_name = $2
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
	err := apps.Pg.QueryRowWArgs(ctx, queryTemplate.RawQuery, ou.OrgID, searchGroupName).Scan(&ts.SearchID, &ts.Query, &ts.MaxResults, &ts.MaxTweetID)
	if err != nil {
		log.Err(err).Interface("sgName", searchGroupName).Msg("SelectTwitterSearchQuery")
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
		tweetID = int(time.Unix(int64(tweetID), 0).UnixNano())
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
