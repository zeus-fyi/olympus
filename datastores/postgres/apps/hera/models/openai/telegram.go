package hera_openai_dbmodels

import (
	"context"
	"encoding/json"
	"unicode/utf8"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func filterSeenTgMsgIds() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "filterSeenTgMsgIds"

	q.RawQuery = `INSERT INTO "public"."ai_incoming_telegram_msgs" ("org_id", "user_id", "timestamp", "chat_id", "message_id", "sender_id", "group_name", "message_text", "metadata")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT ("chat_id", "message_id")
		DO UPDATE SET
			"message_text" = EXCLUDED."message_text"
		RETURNING "telegram_msg_id";`
	return q
}

func InsertNewTgMessages(ctx context.Context, ou org_users.OrgUser, timestamp, chatID, messageID, senderID int, groupName, msgText string, b json.RawMessage) (int, error) {
	q := filterSeenTgMsgIds()
	var msgID int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, timestamp, chatID, messageID, senderID, groupName, msgText, sanitizeJSONRawMessage(b)).Scan(&msgID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msgID, nil
}

func sanitizeJSONRawMessage(raw json.RawMessage) json.RawMessage {
	if utf8.Valid(raw) {
		return raw // The raw message is already valid UTF-8
	}

	buf := make([]byte, 0, len(raw))
	for len(raw) > 0 {
		r, size := utf8.DecodeRune(raw)
		if r == utf8.RuneError && size == 1 {
			// This is an invalid UTF-8 rune, skip it
			raw = raw[size:]
			continue
		}
		buf = append(buf, raw[:size]...)
		raw = raw[size:]
	}
	return buf
}
