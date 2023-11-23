package hera_search

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type DiscordMessage struct {
	// Define the struct fields that match your Discord message attributes
	MessageID int             `json:"message_id"`
	SearchId  int             `json:"search_id"`
	GuildID   string          `json:"guild_id"`
	ChannelID string          `json:"channel_id"`
	Author    json.RawMessage `json:"author"`
	Content   string          `json:"content"`
	Mentions  json.RawMessage `json:"mentions"`
	Reactions json.RawMessage `json:"reactions"`
	Reference json.RawMessage `json:"reference"`
	EditedAt  int             `json:"edited_at"`
	Type      string          `json:"type"`
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
	if err != nil {
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
	if err != nil {
		log.Err(err).Msg("InsertDiscordChannel")
		return err
	}
	return nil
}

func InsertIncomingDiscordMessages(ctx context.Context, messages []*DiscordMessage) ([]int, error) {
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
	for _, message := range messages {
		if message == nil {
			continue
		}
		var messageID int

		err = tx.QueryRow(ctx, q.RawQuery,
			message.MessageID,
			message.SearchId,
			message.GuildID,
			message.ChannelID,
			message.Author,
			message.Content,
			message.Mentions,
			message.Reactions,
			message.Reference,
			message.EditedAt,
			message.Type).Scan(&messageID)
		if err != nil {
			log.Err(err).Msg("InsertIncomingDiscordMessages")
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
