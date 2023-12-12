package hera_search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName, query string, maxResults int) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertRedditSearchQuery"
	q.RawQuery = `INSERT INTO "public"."ai_reddit_search_query" ("org_id", "user_id", "search_group_name", "max_results", "query")
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "search_id";`

	var searchID int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, searchGroupName, maxResults, query).Scan(&searchID)
	if err != nil {
		log.Err(err).Msg("InsertRedditSearchQuery")
		return 0, err
	}
	return searchID, nil
}

func InsertIncomingRedditPosts(ctx context.Context, searchID int, posts []*reddit.Post) ([]string, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertIncomingRedditPosts"

	q.RawQuery = `INSERT INTO "public"."ai_reddit_incoming_posts" ("search_id", "post_id", "post_full_id", "created_at", "edited_at", "permalink", "url", "title", "body", "score", "upvote_ratio", "number_of_comments", "author", "author_id", "reddit_meta", "subreddit")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT ("post_id")
		DO UPDATE SET
			"created_at" = EXCLUDED."created_at",
			"edited_at" = EXCLUDED."edited_at",
			"permalink" = EXCLUDED."permalink",
			"url" = EXCLUDED."url",
			"title" = EXCLUDED."title",
			"body" = EXCLUDED."body",
			"score" = EXCLUDED."score",
			"upvote_ratio" = EXCLUDED."upvote_ratio",
			"number_of_comments" = EXCLUDED."number_of_comments",
			"reddit_meta" = EXCLUDED."reddit_meta"
		RETURNING "post_id";`

	var postIDs []string
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("InsertIncomingRedditPosts")
		return nil, err
	}
	defer tx.Rollback(ctx)
	for _, post := range posts {
		if post == nil {
			continue
		}
		var postID string

		meta := map[string]interface{}{
			"SubredditName":         post.SubredditName,
			"SubredditNamePrefixed": post.SubredditNamePrefixed,
			"SubredditID":           post.SubredditID,
			"SubredditSubscribers":  post.SubredditSubscribers,
			"Spoiler":               post.Spoiler,
			"Locked":                post.Locked,
			"NSFW":                  post.NSFW,
			"IsSelfPost":            post.IsSelfPost,
			"Saved":                 post.Saved,
			"Stickied":              post.Stickied,
		}
		metaJSON, jerr := json.Marshal(meta)
		if jerr != nil {
			return nil, jerr
		}
		var editTsUnix int
		editTs := post.Edited
		if editTs == nil {
			editTs = post.Created
			editTsUnix = int(post.Created.Unix())
		} else {
			editTsUnix = int(post.Edited.Unix())
		}
		err = tx.QueryRow(ctx, q.RawQuery,
			searchID,
			post.ID,
			post.FullID,
			post.Created.Unix(),
			editTsUnix,
			post.Permalink,
			post.URL,
			post.Title,
			post.Body,
			post.Score,
			post.UpvoteRatio,
			post.NumberOfComments,
			post.Author,
			post.AuthorID,
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(metaJSON), Status: IsNull(metaJSON)},
			post.SubredditName,
		).Scan(&postID)
		if err != nil {
			return nil, err
		}
		if postID == "" {
			continue
		}
		postIDs = append(postIDs, postID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("InsertIncomingRedditPosts")
		return nil, err
	}
	return postIDs, nil
}

type RedditSearchQuery struct {
	OrgID         int    `json:"orgID"`
	UserID        int    `json:"userID"`
	PostId        string `json:"postId"`
	LastCreatedAt int    `json:"lastCreatedAt"`
	SearchIndexerParams
}

func SelectRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) ([]*RedditSearchQuery, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectRedditSearchQuery"
	q.RawQuery = `SELECT 
				    sq_sub.search_id AS search_id,
					sq_sub.query as subreddit, 
					sq_sub.last_created_at,
					COALESCE(ip.post_id, '') AS post_id,
					sq_sub.active AS active
				FROM 
					(SELECT
					     sq.search_id,
						 sq.query,
						 sq.active,
						 COALESCE(MAX(ip.created_at), 0) AS last_created_at
					 FROM 
						 public.ai_reddit_search_query sq
					 LEFT JOIN 
						 public.ai_reddit_incoming_posts ip ON sq.query = ip.subreddit
					 WHERE 
						 sq.org_id = $1 AND sq.search_group_name = $2 AND sq.active = true
					 GROUP BY 
						  sq.search_id, sq.query) AS sq_sub
				LEFT JOIN 
					public.ai_reddit_incoming_posts ip ON sq_sub.query = ip.subreddit AND sq_sub.last_created_at = ip.created_at;
				`
	var postId *string
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, searchGroupName)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("SelectRedditSearchQuery: no rows")
		return nil, nil
	}
	var rss []*RedditSearchQuery
	defer rows.Close()
	for rows.Next() {
		rs := &RedditSearchQuery{
			SearchIndexerParams: SearchIndexerParams{
				MaxResults: 100,
			},
		}
		rowErr := rows.Scan(&rs.SearchID, &rs.Query, &rs.LastCreatedAt, &rs.PostId, &rs.Active)
		if rowErr != nil {
			log.Err(rowErr).Msg("Error scanning row in SelectRedditSearchQuery")
			return nil, rowErr
		}
		if postId != nil {
			rs.PostId = *postId
		}
		rss = append(rss, rs)
	}
	if err != nil {
		log.Err(err).Msg("SelectRedditSearchQuery")
		return nil, err
	}
	return rss, nil
}

func redditSearchQuery(ou org_users.OrgUser, sp AiSearchParams) (sql_query_templates.QueryParams, []interface{}) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "redditSearchQuery"
	args := []interface{}{ou.OrgID}

	baseQuery := `SELECT created_at, subreddit, title, body
				  FROM public.ai_reddit_incoming_posts
				  JOIN ai_reddit_search_query sq ON sq.search_id = ai_reddit_incoming_posts.search_id
				  WHERE sq.org_id = $1
				 `
	if sp.Retrieval.RetrievalKeywords != "" {
		args = append(args, sp.Retrieval.RetrievalKeywords)
		baseQuery += fmt.Sprintf(" AND (body_tsvector @@ to_tsquery('english', $%d) OR title_tsvector @@ to_tsquery('english', $%d))", len(args), len(args))
	}

	if !sp.Window.IsWindowEmpty() {
		baseQuery += ` AND`
		tsRangeStart, tsEnd := sp.Window.GetUnixTimestamps()
		baseQuery += fmt.Sprintf(` created_at BETWEEN $%d AND $%d`, len(args)+1, len(args)+2)
		args = append(args, tsRangeStart, tsEnd)
	}
	baseQuery += ` ORDER BY created_at DESC;`
	q.RawQuery = baseQuery
	return q, args
}

func SearchReddit(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q, args := redditSearchQuery(ou, sp)
	var rows pgx.Rows
	var err error
	rows, err = apps.Pg.Query(ctx, q.RawQuery, args...)
	var srs []SearchResult
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SearchReddit")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sr SearchResult
		sr.Source = "reddit"
		title := ""
		body := ""
		rowErr := rows.Scan(
			&sr.UnixTimestamp, &sr.Group, &title, &body,
		)
		if len(body) <= 0 {
			continue
		}
		sr.Value = title + "\n " + body + "\n"
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SearchReddit"))
			return nil, rowErr
		}
		srs = append(srs, sr)
	}
	return srs, nil
}
