package hera_search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type DiscordMessage struct {
	// Define the struct fields that match your Discord message attributes
	MessageID int    `json:"message_id"`
	SearchId  int    `json:"search_id"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	Author    any    `json:"author"`
	Content   string `json:"content"`
	Mentions  any    `json:"mentions"`
	Reactions any    `json:"reactions"`
	Reference any    `json:"reference"`
	EditedAt  int    `json:"edited_at"`
	Type      string `json:"type"`
}

type DiscordMetadata struct {
	GuildName    string `json:"guildName"`
	Category     string `json:"topic"`
	CategoryName string `json:"categoryName"`
}

// discordSearchQuery will be extended to support new search parameters when they are not empty
func discordSearchQuery(ou org_users.OrgUser, sp AiSearchParams) (sql_query_templates.QueryParams, []interface{}) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "discordSearchQuery"
	args := []interface{}{ou.OrgID}

	baseQuery := `SELECT cm.timestamp_creation, cm.content, gi.name, ci.category, ci.name
				  FROM public.ai_incoming_discord_messages cm
				  JOIN public.ai_discord_channel ci ON ci.channel_id = cm.channel_id
				  JOIN public.ai_discord_guild gi ON gi.guild_id = cm.guild_id
				  JOIN public.ai_discord_search_query sq ON sq.search_id = cm.search_id
				  WHERE sq.org_id = $1`

	// Handle positive keywords
	if sp.Retrieval.RetrievalKeywords != nil && *sp.Retrieval.RetrievalKeywords != "" {
		posQuery := formatKeywordsForTsQuery(aws.StringValue(sp.Retrieval.RetrievalKeywords), false)
		if posQuery != "" {
			baseQuery += fmt.Sprintf(` AND cm.content_tsvector @@ to_tsquery('english', $%d)`, len(args)+1)
			args = append(args, posQuery)
		}
	}

	// Handle negative keywords
	if sp.Retrieval.RetrievalNegativeKeywords != nil && *sp.Retrieval.RetrievalNegativeKeywords != "" {
		negQuery := formatKeywordsForTsQuery(aws.StringValue(sp.Retrieval.RetrievalNegativeKeywords), true)
		if negQuery != "" {
			baseQuery += fmt.Sprintf(` AND cm.content_tsvector @@ to_tsquery('english', $%d)`, len(args)+1)
			args = append(args, negQuery)
		}
	}

	if !sp.Window.IsWindowEmpty() {
		tsRangeStart, tsEnd := sp.Window.GetUnixTimestamps()
		baseQuery += fmt.Sprintf(` AND cm.timestamp_creation BETWEEN $%d AND $%d`, len(args)+1, len(args)+2)
		args = append(args, tsRangeStart, tsEnd)
	}

	if sp.Retrieval.RetrievalPlatformGroups != nil && *sp.Retrieval.RetrievalPlatformGroups != "" {
		groupFilters := strings.Split(*sp.Retrieval.RetrievalPlatformGroups, ",")
		baseQuery += ` AND (`
		queryParts := make([]string, 0, len(groupFilters))
		for _, filter := range groupFilters {
			trimmedFilter := strings.TrimSpace(filter)
			if trimmedFilter != "" {
				argCount := len(args) + 1
				queryPart := fmt.Sprintf(`gi.name ILIKE $%d`, argCount)
				queryParts = append(queryParts, queryPart)
				args = append(args, "%"+trimmedFilter+"%")
			}
		}
		baseQuery += strings.Join(queryParts, " OR ") + `)`
	}

	if sp.Retrieval.DiscordFilters != nil && *sp.Retrieval.DiscordFilters.CategoryName != "" {
		categoryNames := strings.Split(*sp.Retrieval.DiscordFilters.CategoryName, ",")
		baseQuery += ` AND (`
		queryParts := make([]string, 0, len(categoryNames))
		for _, filter := range categoryNames {
			trimmedFilter := strings.TrimSpace(filter)
			if trimmedFilter != "" {
				argCount := len(args) + 1
				queryPart := fmt.Sprintf(`ci.name ILIKE $%d`, argCount)
				queryParts = append(queryParts, queryPart)
				args = append(args, "%"+trimmedFilter+"%")
			}
		}
		baseQuery += strings.Join(queryParts, " OR ") + `)`
	}

	baseQuery += ` ORDER BY cm.timestamp_creation DESC;`
	q.RawQuery = baseQuery
	return q, args
}

func SearchDiscord(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q, args := discordSearchQuery(ou, sp)
	var srs []SearchResult
	var rows pgx.Rows
	var err error

	rows, err = apps.Pg.Query(ctx, q.RawQuery, args...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SearchDiscord")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sr := SearchResult{Source: "discord", Verified: aws.Bool(true)}
		rowErr := rows.Scan(&sr.UnixTimestamp, &sr.Value, &sr.Group, &sr.DiscordMetadata.Category, &sr.DiscordMetadata.CategoryName)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SearchDiscord"))
			return nil, rowErr
		}
		srs = append(srs, sr)
	}
	return srs, nil
}

func InsertDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName, query string, maxResults int) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertDiscordSearchQuery"
	q.RawQuery = `INSERT INTO "public"."ai_discord_search_query" ("org_id", "user_id", "search_group_name", "max_results", "query")
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "search_id";`

	var searchID int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, searchGroupName, maxResults, query).Scan(&searchID)
	if err != nil {
		log.Err(err).Msg("InsertDiscordSearchQuery")
		return 0, err
	}
	return searchID, nil
}

func InsertDiscordGuild(ctx context.Context, guildID, name string) error {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertDiscordGuild"
	q.RawQuery = `INSERT INTO "public"."ai_discord_guild" ("guild_id", "name")
                  VALUES ($1, $2)
                  ON CONFLICT ("guild_id")
                  DO UPDATE SET
                      "name" = EXCLUDED.name;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, guildID, name)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("InsertDiscordGuild")
		return err
	}
	return nil
}

func InsertDiscordChannel(ctx context.Context, searchID int, guildID, channelID, categoryID, category, name, topic string) error {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertDiscordChannel"
	q.RawQuery = `INSERT INTO "public"."ai_discord_channel" ("search_id", "guild_id", "channel_id", "category_id", "category", "name", "topic")
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT ("channel_id")
        DO UPDATE SET
            "name" = EXCLUDED.name,
            "topic" = EXCLUDED.topic,
            "category" = EXCLUDED.category;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, searchID, guildID, channelID, categoryID, category, name, topic)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("InsertDiscordChannel")
		return err
	}
	return nil
}

func InsertIncomingDiscordMessages(ctx context.Context, searchID int, messages hera_discord.ChannelMessages) ([]int, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "insertIncomingDiscordMessages"
	q.RawQuery = `INSERT INTO "public"."ai_incoming_discord_messages" ("timestamp_creation","message_id", "search_id", "guild_id", "channel_id", "author", "content", "mentions", "reactions", "reference", "timestamp_edited", "type")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT ("message_id")
		DO UPDATE SET
			"content" = EXCLUDED.content,
			"mentions" = EXCLUDED.mentions,
			"reactions" = EXCLUDED.reactions,
			"reference" = EXCLUDED.reference,
			"timestamp_edited" = EXCLUDED.timestamp_edited
		RETURNING "timestamp_creation";`

	var messageIDs []int
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	for _, message := range messages.Messages {

		ba, berr := json.Marshal(message.Author)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
		br, berr := json.Marshal(message.Reactions)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
		if len(message.Reactions) <= 0 {
			br = nil
		}
		bm, berr := json.Marshal(message.Mentions)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
		if len(message.Mentions) <= 0 {
			bm = nil
		}
		rrf, berr := json.Marshal(message.Reference)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
		if message.Reference.IsEmpty() {
			rrf = nil
		}
		if len(message.Content) <= 0 {
			continue
		}
		var messageID int
		err = tx.QueryRow(ctx, q.RawQuery,
			int(message.TimestampCreated.Unix()),
			message.Id,
			searchID,
			messages.Guild.Id,
			messages.Channel.Id,
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(ba), Status: IsNull(ba)},
			message.Content,
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(bm), Status: IsNull(bm)},
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(br), Status: IsNull(br)},
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(rrf), Status: IsNull(rrf)},
			int(message.TimestampEdited.Unix()),
			message.Type).Scan(&messageID)
		if err != nil {
			log.Err(err).Interface("message", message).Msg("InsertIncomingDiscordMessages")
			return nil, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return messageIDs, nil
}

func IsNull(b []byte) pgtype.Status {
	if b == nil {
		return pgtype.Null
	}
	return pgtype.Present
}

type DiscordSearchResult struct {
	SearchID        int    `json:"search_id"`
	GuildID         string `json:"guild_id"`
	ChannelID       string `json:"channel_id"`
	MaxMessageID    int    `json:"max_message_id"`
	SearchGroupName string `json:"searchGroupName,omitempty"`
}

type DiscordSearchResultWrapper struct {
	SearchID int                    `json:"search_id"`
	Results  []*DiscordSearchResult `json:"results"`
}

func SelectDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) (*DiscordSearchResultWrapper, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectDiscordSearchQuery"
	q.RawQuery = `
        SELECT dm.search_id, dm.guild_id, dm.channel_id, MAX(dm.timestamp_creation) AS max_message_id
        FROM public.ai_incoming_discord_messages dm
        INNER JOIN public.ai_discord_search_query dsq 
        ON dm.search_id = dsq.search_id
        WHERE dsq.org_id = $1 AND dsq.search_group_name = $2
        GROUP BY dm.search_id, dm.guild_id, dm.channel_id;
    `

	var results []*DiscordSearchResult
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, searchGroupName)
	if err != nil {
		log.Err(err).Msg("SelectDiscordSearchQuery")
		return nil, err
	}
	defer rows.Close()
	var searchID int
	for rows.Next() {
		var r DiscordSearchResult
		if err = rows.Scan(&r.SearchID, &r.GuildID, &r.ChannelID, &r.MaxMessageID); err != nil {
			log.Err(err).Msg("Error scanning row in SelectDiscordSearchQuery")
			return nil, err
		}
		if searchID == 0 {
			searchID = r.SearchID
		}
		results = append(results, &r)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("Error iterating rows in SelectDiscordSearchQuery")
		return nil, err
	}

	return &DiscordSearchResultWrapper{
		SearchID: searchID,
		Results:  results,
	}, nil
}

func SelectDiscordSearchQueryByGuildChannel(ctx context.Context, ou org_users.OrgUser, guildID, channelID string) (*DiscordSearchResultWrapper, error) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "selectDiscordSearchQuery"
	q.RawQuery = `
        SELECT dsq.search_group_name, dsq.search_id, gi.guild_id, cm.channel_id
		FROM public.ai_incoming_discord_messages cm
		JOIN public.ai_discord_channel ci ON ci.channel_id = cm.channel_id
		JOIN public.ai_discord_guild gi ON gi.guild_id = cm.guild_id
		JOIN public.ai_discord_search_query dsq ON dsq.search_id = cm.search_id 
        WHERE dsq.org_id = $1 AND gi.guild_id = $2 AND ci.channel_id = $3;
    `

	var results []*DiscordSearchResult
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, guildID, channelID)
	if err != nil {
		log.Err(err).Msg("SelectDiscordSearchQuery")
		return nil, err
	}
	defer rows.Close()
	var searchID int
	for rows.Next() {
		var r DiscordSearchResult
		if err = rows.Scan(&r.SearchGroupName, &r.SearchID, &r.GuildID, &r.ChannelID); err != nil {
			log.Err(err).Msg("Error scanning row in SelectDiscordSearchQuery")
			return nil, err
		}
		if searchID == 0 {
			searchID = r.SearchID
		}
		results = append(results, &r)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("Error iterating rows in SelectDiscordSearchQuery")
		return nil, err
	}

	return &DiscordSearchResultWrapper{
		SearchID: searchID,
		Results:  results,
	}, nil
}
