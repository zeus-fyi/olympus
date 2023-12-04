package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type RetrievalItem struct {
	RetrievalID       int    `json:"retrievalID"`    // ID of the retrieval
	RetrievalName     string `json:"retrievalName"`  // Name of the retrieval
	RetrievalGroup    string `json:"retrievalGroup"` // Group of the retrieval
	RetrievalPlatform string `json:"retrievalPlatform"`
	Instructions      []byte `json:"instructions"` // Instructions for the retrieval
}

func InsertRetrieval(ctx context.Context, ou org_users.OrgUser, item *RetrievalItem, b []byte) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO public.ai_retrieval_library (org_id, user_id, retrieval_name, retrieval_group, retrieval_platform, instructions)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, retrieval_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            retrieval_group = EXCLUDED.retrieval_group,
            retrieval_platform = EXCLUDED.retrieval_platform,
            instructions = EXCLUDED.instructions
        RETURNING retrieval_id;`
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, item.RetrievalName, item.RetrievalGroup, item.RetrievalPlatform, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&item.RetrievalID)
	if err != nil {
		log.Err(err).Msg("failed to insert retrieval")
		return err
	}
	return nil
}

func SelectRetrievals(ctx context.Context, ou org_users.OrgUser) ([]RetrievalItem, error) {
	query := `
        SELECT retrieval_id, retrieval_name, retrieval_group, retrieval_platform, instructions
        FROM public.ai_retrieval_library
        WHERE org_id = $1
        ORDER BY retrieval_id DESC;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, query, ou.OrgID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var retrievals []RetrievalItem

	// Iterating over the result set
	for rows.Next() {
		var retrieval RetrievalItem
		var instructions pgtype.JSONB
		err = rows.Scan(&retrieval.RetrievalID, &retrieval.RetrievalName, &retrieval.RetrievalGroup, &retrieval.RetrievalPlatform, &instructions)
		if err != nil {
			log.Err(err).Msg("failed to scan retrieval")
			return nil, err
		}
		retrieval.Instructions = instructions.Bytes // Assuming Instructions field in RetrievalItem is of type []byte
		retrievals = append(retrievals, retrieval)
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating retrieval rows")
		return nil, err
	}

	return retrievals, nil
}
