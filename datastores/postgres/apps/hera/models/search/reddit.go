package hera_search

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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

	q.RawQuery = `INSERT INTO "public"."ai_reddit_incoming_posts" ("search_id", "post_id", "post_full_id", "created_at", "edited_at", "permalink", "url", "title", "body", "score", "upvote_ratio", "number_of_comments", "author", "author_id", "reddit_meta")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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
			metaJSON,
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
		return nil, err
	}
	return postIDs, nil
}

type RedditSearchQuery struct {
	SearchID        int    `json:"searchID"`
	OrgID           int    `json:"orgID"`
	UserID          int    `json:"userID"`
	SearchGroupName string `json:"searchGroupName"`
	MaxResults      int    `json:"maxResults"`
	LastCreatedAt   int    `json:"lastCreatedAt"`
	FullPostId      string `json:"fullPostId"`
	Query           string `json:"query"`
}

func SelectRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) (*RedditSearchQuery, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectRedditSearchQuery"
	q.RawQuery = `
        SELECT sq.search_id, sq.query, sq.max_results, ip.post_full_id, COALESCE(MAX(ip.created_at), 0) AS last_created_at
        FROM public.ai_reddit_search_query sq
        LEFT JOIN public.ai_reddit_incoming_posts ip ON sq.search_id = ip.search_id
        WHERE sq.org_id = $1 AND sq.user_id = $2 AND sq.search_group_name = $3
        GROUP BY sq.search_id, sq.query, sq.max_results, ip.post_full_id;
    `
	rs := &RedditSearchQuery{}

	var postId *string
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, searchGroupName).Scan(&rs.SearchID, &rs.Query, &rs.MaxResults, &postId, &rs.LastCreatedAt)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("SelectRedditSearchQuery: no rows")
		return nil, nil
	}
	if postId != nil {
		rs.FullPostId = *postId
	}
	if err != nil {
		log.Err(err).Msg("SelectRedditSearchQuery")
		return nil, err
	}
	return rs, err
}

type DiscordSearchResult struct {
	SearchID     int    `json:"search_id"`
	GuildID      string `json:"guild_id"`
	ChannelID    string `json:"channel_id"`
	MaxMessageID int    `json:"max_message_id"`
}

func SelectDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) ([]*DiscordSearchResult, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectDiscordSearchQuery"
	q.RawQuery = `
        SELECT dm.search_id, dm.guild_id, dm.channel_id, MAX(dm.message_id) AS max_message_id
        FROM public.ai_incoming_discord_messages dm
        INNER JOIN public.ai_discord_search_query dsq 
        ON dm.search_id = dsq.search_id
        WHERE dsq.org_id = $1 AND dsq.user_id = $2 AND dsq.search_group_name = $3
        GROUP BY dm.search_id, dm.guild_id, dm.channel_id;
    `

	var results []*DiscordSearchResult
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, ou.UserID, searchGroupName)
	if err != nil {
		log.Err(err).Msg("SelectDiscordSearchQuery")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r DiscordSearchResult
		if err = rows.Scan(&r.SearchID, &r.GuildID, &r.ChannelID, &r.MaxMessageID); err != nil {
			log.Err(err).Msg("Error scanning row in SelectDiscordSearchQuery")
			return nil, err
		}
		results = append(results, &r)
	}

	if err = rows.Err(); err != nil {
		log.Err(err).Msg("Error iterating rows in SelectDiscordSearchQuery")
		return nil, err
	}

	return results, nil
}
