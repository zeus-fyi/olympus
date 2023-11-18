package hera_openai_dbmodels

import (
	"context"

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

func InsertNewTgMessages(ctx context.Context, ou org_users.OrgUser, timestamp, chatID, messageID, senderID int, groupName, msgText string, b []byte) (int, error) {
	q := filterSeenTgMsgIds()
	var msgID int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, timestamp, chatID, messageID, senderID, groupName, msgText, b).Scan(&msgID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msgID, nil
}
