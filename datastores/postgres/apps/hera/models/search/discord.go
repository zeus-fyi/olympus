package hera_search

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
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

func InsertDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string, maxResults int, query string) (int, error) {
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
	q.RawQuery = `INSERT INTO "public"."ai_incoming_discord_messages" ("message_id", "search_id", "guild_id", "channel_id", "author", "content", "mentions", "reactions", "reference", "timestamp_edited", "type")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT ("message_id")
		DO UPDATE SET
			"content" = EXCLUDED.content,
			"mentions" = EXCLUDED.mentions,
			"reactions" = EXCLUDED.reactions,
			"reference" = EXCLUDED.reference,
			"timestamp_edited" = EXCLUDED.timestamp_edited
		RETURNING "message_id";`

	var messageIDs []int
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	for _, message := range messages.Messages {
		var messageID int
		mi, berr := strconv.Atoi(message.Id)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
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
		bm, berr := json.Marshal(message.Mentions)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}
		rrf, berr := json.Marshal(message.Reference)
		if berr != nil {
			log.Err(berr).Msg("InsertIncomingDiscordDataFromSearch")
			return nil, berr
		}

		err = tx.QueryRow(ctx, q.RawQuery,
			mi,
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
	SearchID     int    `json:"search_id"`
	GuildID      string `json:"guild_id"`
	ChannelID    string `json:"channel_id"`
	MaxMessageID int    `json:"max_message_id"`
}

type DiscordSearchResultWrapper struct {
	SearchID int                    `json:"search_id"`
	Results  []*DiscordSearchResult `json:"results"`
}

func SelectDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, searchGroupName string) (*DiscordSearchResultWrapper, error) {
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
