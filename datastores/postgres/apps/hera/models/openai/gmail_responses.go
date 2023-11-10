package hera_openai_dbmodels

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func filterSeenIds() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "filterSeenIds"
	q.RawQuery = `INSERT INTO "public"."ai_incoming_email_tasks" ("msg_id", "from", "subject", "contents")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ("msg_id") DO NOTHING
		RETURNING "email_id"`
	return q
}

func InsertNewEmails(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	q := filterSeenIds()
	var emailID int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, msg.MsgId, msg.From, msg.Subject, msg.Body).Scan(&emailID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return emailID, nil
}

/*
CREATE TABLE "public"."ai_outgoing_email_tasks" (
    "response_id" int8 NOT NULL DEFAULT next_id(),
    "email_id" int8 NOT NULL,
    "msg_id" int8 NOT NULL,
    "from" text NOT NULL,
    "subject" text NOT NULL DEFAULT 'empty',
    "contents" text NOT NULL DEFAULT 'empty'
);
*/
//func insertOutgoing() sql_query_templates.QueryParams {
//	q := sql_query_templates.QueryParams{}
//	q.QueryName = "insertOutgoing"
//	q.RawQuery = `INSERT INTO "public"."ai_outgoing_email_tasks" ("msg_id", "email_id", "from", "subject", "contents")
//		VALUES ($1, $2, $3, $4)
//		ON CONFLICT ("msg_id") DO NOTHING
//		RETURNING "email_id"`
//	return q
//}
//
//func InsertNewEmails(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
//	q := filterSeenIds()
//	var emailID int
//	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, msg.MsgId, msg.From, msg.Subject, msg.Body).Scan(&emailID)
//	if err == pgx.ErrNoRows {
//		return 0, nil
//	}
//	if err != nil {
//		return 0, err
//	}
//	return emailID, nil
//}
