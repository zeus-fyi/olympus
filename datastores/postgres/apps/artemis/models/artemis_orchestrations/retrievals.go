package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type RetrievalItem struct {
	RetrievalStrID           *string                           `json:"retrievalStrID"`
	RetrievalID              *int                              `json:"retrievalID,omitempty"` // ID of the retrieval
	RetrievalName            string                            `json:"retrievalName"`         // Name of the retrieval
	RetrievalGroup           string                            `json:"retrievalGroup"`        // Group of the retrieval
	RetrievalItemInstruction `json:"retrievalItemInstruction"` // Instructions for the retrieval
}

type RetrievalItemInstruction struct {
	RetrievalPlatform         string          `json:"retrievalPlatform"`
	RetrievalPrompt           *string         `json:"retrievalPrompt,omitempty"`           // Prompt for the retrieval
	RetrievalPlatformGroups   *string         `json:"retrievalPlatformGroups,omitempty"`   // Platform groups for the retrieval
	RetrievalKeywords         *string         `json:"retrievalKeywords,omitempty"`         // Keywords for the retrieval
	RetrievalNegativeKeywords *string         `json:"retrievalNegativeKeywords,omitempty"` // Keywords for the retrieval
	RetrievalUsernames        *string         `json:"retrievalUsernames,omitempty"`        // Usernames for the retrieval
	DiscordFilters            *DiscordFilters `json:"discordFilters,omitempty"`            // Discord filters for the retrieval
	WebFilters                *WebFilters     `json:"webFilters,omitempty"`                // Web filters for the retrieval

	Instructions json.RawMessage `json:"instructions,omitempty"` // Instructions for the retrieval
}

type WebFilters struct {
	RoutingGroup         *string  `json:"routingGroup,omitempty"`
	LbStrategy           *string  `json:"lbStrategy,omitempty"`
	MaxRetries           *int     `json:"maxRetries,omitempty"`
	BackoffCoefficient   *float64 `json:"backoffCoefficient,omitempty"`
	EndpointRoutePath    *string  `json:"endpointRoutePath,omitempty"`
	EndpointREST         *string  `json:"endpointREST,omitempty"`
	PayloadPreProcessing *string  `json:"payloadPreProcessing,omitempty"`
}

type DiscordFilters struct {
	CategoryTopic *string `json:"categoryTopic,omitempty"`
	CategoryName  *string `json:"categoryName,omitempty"`
	Category      *string `json:"category,omitempty"`
}

func SetInstructions(r *RetrievalItem) error {
	b, err := json.Marshal(r.RetrievalItemInstruction)
	if err != nil {
		log.Err(err).Msg("failed to marshal retrieval instructions")
		return err
	}
	r.Instructions = b
	return nil
}

func InsertRetrieval(ctx context.Context, ou org_users.OrgUser, item *RetrievalItem) error {
	q := sql_query_templates.QueryParams{}
	err := SetInstructions(item)
	if err != nil {
		log.Err(err).Msg("failed to set retrieval instructions")
		return err
	}
	if item.Instructions == nil {
		return errors.New("instructions cannot be nil")
	}
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
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, item.RetrievalName, item.RetrievalGroup, item.RetrievalPlatform,
		&pgtype.JSONB{Bytes: sanitizeBytesUTF8(item.Instructions), Status: IsNull(item.Instructions)}).Scan(&item.RetrievalID)
	if err != nil {
		log.Err(err).Msg("failed to insert retrieval")
		return err
	}
	return nil
}

func SelectRetrievals(ctx context.Context, ou org_users.OrgUser, retID int) ([]RetrievalItem, error) {
	args := []interface{}{ou.OrgID}
	var queryAddOn string
	if retID > 0 {
		args = append(args, retID)
		queryAddOn = " AND retrieval_id = $2"
	}
	query := `
        SELECT retrieval_id, retrieval_id::text AS retrieval_id_str, retrieval_name, retrieval_group, retrieval_platform, instructions
        FROM public.ai_retrieval_library
        WHERE org_id = $1 ` + queryAddOn
	// Executing the query
	rows, err := apps.Pg.Query(ctx, query, args...)
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
		err = rows.Scan(&retrieval.RetrievalID, &retrieval.RetrievalStrID, &retrieval.RetrievalName,
			&retrieval.RetrievalGroup, &retrieval.RetrievalPlatform, &retrieval.Instructions)
		if err != nil {
			log.Err(err).Msg("failed to scan retrieval")
			return nil, err
		}
		if retrieval.RetrievalID != nil {
			retrieval.RetrievalStrID = aws.String(fmt.Sprintf("%d", *retrieval.RetrievalID))
		}
		if retrieval.Instructions != nil {
			b := retrieval.Instructions
			if b != nil {
				err = json.Unmarshal(b, &retrieval.RetrievalItemInstruction)
				if err != nil {
					log.Err(err).Msg("failed to unmarshal retrieval instructions")
					return nil, err
				}
			}
			retrieval.Instructions = nil
		}
		retrievals = append(retrievals, retrieval)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating retrieval rows")
		return nil, err
	}
	return retrievals, nil
}
