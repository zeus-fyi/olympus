package hera_openai_dbmodels

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func insertCompletionResp() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertCompletionResponse"
	q.RawQuery =
		`WITH cte_insert_token_usage AS (
			SELECT tokens_remaining, tokens_consumed FROM hera_openai_usage WHERE org_id = $1
		), cte_update_token_usage AS (
     		INSERT INTO hera_openai_usage (org_id, tokens_remaining, tokens_consumed)
            VALUES ($1, 0, 0)
            ON CONFLICT (org_id) 
            DO UPDATE SET
                tokens_remaining = hera_openai_usage.tokens_remaining - $5,
                tokens_consumed = hera_openai_usage.tokens_consumed + $5
            WHERE hera_openai_usage.org_id = $1
            RETURNING tokens_remaining, tokens_consumed
		)
		INSERT INTO completion_responses(org_id, user_id, prompt_tokens, completion_tokens, total_tokens, model, completion_choices)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING response_id;
	`
	return q
}

const Sn = "OpenAI"

func IsNull(b []byte) pgtype.Status {
	if b == nil {
		return pgtype.Null
	}
	return pgtype.Present
}

func sanitizeBytesUTF8(b []byte) []byte {
	bs := bytes.ReplaceAll(b, []byte{0}, []byte{})
	return bs
}

func InsertCompletionResponseChatGpt(ctx context.Context, ou org_users.OrgUser, response openai.ChatCompletionResponse) (int, error) {
	q := insertCompletionResp()
	completionChoices, err := json.Marshal(response.Choices)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return 0, err
	}
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var rid int
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens, response.Model, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(completionChoices), Status: IsNull(completionChoices)}).Scan(&rid)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return rid, err
	}
	return rid, err
}

func InsertCompletionResponse(ctx context.Context, ou org_users.OrgUser, response openai.CompletionResponse) error {
	q := insertCompletionResp()
	completionChoices, err := json.Marshal(response.Choices)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))

	r, err := apps.Pg.Exec(ctx, q.RawQuery, ou.OrgID, ou.UserID, response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens, response.Model, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(completionChoices), Status: IsNull(completionChoices)})
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("OrgUser: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return err
}

func checkTokenBalance() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertCompletionResponse"
	q.RawQuery = `SELECT tokens_remaining, tokens_consumed FROM hera_openai_usage WHERE org_id = $1`
	return q
}

func CheckTokenBalance(ctx context.Context, ou org_users.OrgUser) (autogen_bases.HeraOpenaiUsage, error) {
	q := checkTokenBalance()
	usageOpenAI := autogen_bases.HeraOpenaiUsage{
		OrgID:           ou.OrgID,
		TokensRemaining: 0,
		TokensConsumed:  0,
	}
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID).Scan(&usageOpenAI.TokensRemaining, &usageOpenAI.TokensConsumed)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return usageOpenAI, err
	}
	return usageOpenAI, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
