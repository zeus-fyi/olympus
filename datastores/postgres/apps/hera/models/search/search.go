package hera_search

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AiSearchParams struct {
	SearchContentText    string `json:"searchContentText,omitempty"`
	GroupFilter          string `json:"groupFilter,omitempty"`
	Usernames            string `json:"usernames,omitempty"`
	WorkflowInstructions string `json:"workflowInstructions,omitempty"`
}

type SearchResult struct {
	UnixTimestamp int                                   `json:"unixTimestamp"`
	Source        string                                `json:"source"`
	Value         string                                `json:"value"`
	Group         string                                `json:"group"`
	Metadata      hera_openai_dbmodels.TelegramMetadata `json:"metadata"`
}

//	func telegramSearchQuery() sql_query_templates.QueryParams {
//		q := sql_query_templates.QueryParams{}
//		q.QueryName = "telegramSearchQuery"
//		q.RawQuery = `SELECT group_name, message_text
//					  FROM public.ai_incoming_telegram_msgs
//					  WHERE org_id = $1 AND group_name ILIKE '%$1%'
//					  AND to_tsvector('english', message_text) @@ plainto_tsquery('english', '$3')
//					  ORDER BY chat_id, message_id DESC;`
//		return q
//	}
func telegramSearchQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "telegramSearchQuery"
	q.RawQuery = `SELECT timestamp, group_name, message_text, metadata
				  FROM public.ai_incoming_telegram_msgs
              	  WHERE org_id = $1 AND group_name ILIKE '%' || $2 || '%' 
				  ORDER BY chat_id, message_id DESC;`
	return q
}

func SearchTelegram(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q := telegramSearchQuery()
	var srs []SearchResult
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, sp.GroupFilter)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SearchTelegram")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sr SearchResult
		sr.Source = "telegram"
		rowErr := rows.Scan(
			&sr.UnixTimestamp, &sr.Group, &sr.Value, &sr.Metadata,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SearchTelegram"))
			return nil, rowErr
		}
		srs = append(srs, sr)
	}
	return srs, nil
}
