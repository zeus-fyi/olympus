package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type RetrievalItem struct {
	RetrievalID    int    `json:"retrievalID"`    // ID of the retrieval
	RetrievalName  string `json:"retrievalName"`  // Name of the retrieval
	RetrievalGroup string `json:"retrievalGroup"` // Group of the retrieval
}

func InsertRetrieval(ctx context.Context, ou org_users.OrgUser, item *RetrievalItem, b []byte) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO public.ai_retrieval_library (org_id, user_id, retrieval_name, retrieval_group, instructions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (org_id, retrieval_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            retrieval_group = EXCLUDED.retrieval_group,
            instructions = EXCLUDED.instructions
        RETURNING retrieval_id;`
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, item.RetrievalName, item.RetrievalGroup, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&item.RetrievalID)
	if err != nil {
		log.Err(err).Msg("failed to insert retrieval")
		return err
	}
	return nil
}
