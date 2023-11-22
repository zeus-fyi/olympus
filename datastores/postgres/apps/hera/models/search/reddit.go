package hera_search

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func insertRedditSearchQuery() (sql_query_templates.QueryParams, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertRedditSearchQuery"
	q.RawQuery = `INSERT INTO "public"."ai_reddit_search_query" ("org_id", "user_id", "search_group_name", "max_results", "query")
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "search_id";`
	return q, nil
}

func InsertRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName, query string, maxResults int) (int, error) {
	queryTemplate, err := insertRedditSearchQuery()
	if err != nil {
		return 0, err
	}
	var searchID int
	err = apps.Pg.QueryRowWArgs(ctx, queryTemplate.RawQuery, ou.OrgID, ou.UserID, searchGroupName, maxResults, query).Scan(&searchID)
	if err != nil {
		log.Err(err).Msg("InsertRedditSearchQuery")
		return 0, err
	}
	return searchID, nil
}

func insertIncomingRedditPosts() (sql_query_templates.QueryParams, error) {
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
	return q, nil
}

func InsertIncomingRedditPosts(ctx context.Context, searchID int, posts []*reddit.Post) ([]string, error) {
	queryTemplate, err := insertIncomingRedditPosts()
	if err != nil {
		return nil, err
	}
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
		err = tx.QueryRow(ctx, queryTemplate.RawQuery,
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
