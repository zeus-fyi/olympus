package hera_openai_dbmodels

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgtype"
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

type TelegramMessage struct {
	Timestamp   int    `json:"timestamp"`
	GroupName   string `json:"group_name"`
	SenderID    int    `json:"sender_id"`
	MessageText string `json:"message_text"`
	ChatID      int    `json:"chat_id"`
	MessageID   int    `json:"message_id"`
	TelegramMetadata
}

type TelegramMetadata struct {
	IsReply       bool   `json:"is_reply,omitempty"`
	IsChannel     bool   `json:"is_channel,omitempty"`
	IsGroup       bool   `json:"is_group,omitempty"`
	IsPrivate     bool   `json:"is_private,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	MutualContact bool   `json:"mutual_contact,omitempty"`
	Username      string `json:"username,omitempty"`
}

func (t *TelegramMetadata) Sanitize() {
	t.FirstName = sanitizeUTF8(t.FirstName)
	t.LastName = sanitizeUTF8(t.LastName)
	t.Phone = sanitizeUTF8(t.Phone)
	t.Username = sanitizeUTF8(t.Username)
}
func InsertNewTgMessages(ctx context.Context, ou org_users.OrgUser, msg TelegramMessage) (int, error) {
	q := filterSeenTgMsgIds()
	var msgID int
	msg.TelegramMetadata.Sanitize()
	msg.GroupName = sanitizeUTF8(msg.GroupName)
	msg.MessageText = sanitizeUTF8(msg.MessageText)

	metadataJSON, err := json.Marshal(msg.TelegramMetadata)
	if err != nil {
		// handle error, maybe return it
		return 0, err
	}
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, msg.Timestamp, msg.ChatID, msg.MessageID, msg.SenderID, msg.GroupName, msg.MessageText, pgtype.JSONB{Bytes: metadataJSON, Status: pgtype.Present}).Scan(&msgID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return msgID, nil
}
func sanitizeUTF8(s string) string {
	bs := bytes.ReplaceAll([]byte(s), []byte{0}, []byte{})
	return strings.ToValidUTF8(string(bs), "?")
}
