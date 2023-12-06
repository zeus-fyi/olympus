package hera_search

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"golang.org/x/crypto/sha3"
)

type AiSearchParams struct {
	SearchContentText    string       `json:"searchContentText,omitempty"`
	WorkflowInstructions string       `json:"workflowInstructions,omitempty"`
	TimeRange            string       `json:"timeRange,omitempty"`
	GroupFilter          string       `json:"groupFilter,omitempty"`
	Platforms            string       `json:"platforms,omitempty"`
	Usernames            string       `json:"usernames,omitempty"`
	SearchInterval       TimeInterval `json:"searchInterval,omitempty"`
	AnalysisInterval     TimeInterval `json:"analysisInterval,omitempty"`
}
type AiModelParams struct {
	Model         string `json:"model"`
	TokenCountMax int    `json:"tokenCountMax"`
}
type TimeInterval [2]time.Time

func (ti *TimeInterval) GetUnixTimestamps() (int, int) {
	if ti == nil {
		return 0, 0
	}
	return int(ti[0].Unix()), int(ti[1].Unix())
}

type SearchResult struct {
	UnixTimestamp   int              `json:"unixTimestamp"`
	Source          string           `json:"source"`
	Value           string           `json:"value"`
	Group           string           `json:"group"`
	Metadata        TelegramMetadata `json:"metadata,omitempty"`
	DiscordMetadata DiscordMetadata  `json:"discordMetadata"`
}
type ByTimestamp []SearchResult

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].UnixTimestamp > a[j].UnixTimestamp }

// SortSearchResults sorts the slice of SearchResult in descending order by UnixTimestamp.
func SortSearchResults(results []SearchResult) {
	sort.Sort(ByTimestamp(results))
}

type SearchResults struct {
	AiModelParams AiModelParams  `json:"aiModelParams"`
	Results       []SearchResult `json:"results"`
}

func FormatSearchResultsV2(results []SearchResult) string {
	var builder strings.Builder

	for _, result := range results {
		var parts []string

		// Always include the UnixTimestamp
		parts = append(parts, fmt.Sprintf("%d", result.UnixTimestamp))

		// Conditionally append other fields if they are not empty
		if result.Source != "" {
			parts = append(parts, escapeString(result.Source))
		}
		if result.Group != "" {
			parts = append(parts, escapeString(result.Group))
		}
		if result.DiscordMetadata.Category != "" {
			parts = append(parts, escapeString(result.DiscordMetadata.Category))
		}
		if result.DiscordMetadata.CategoryName != "" {
			parts = append(parts, escapeString(result.DiscordMetadata.CategoryName))
		}
		if result.Metadata.Username != "" {
			parts = append(parts, escapeString(result.Metadata.Username))
		}
		if result.Value != "" {
			parts = append(parts, escapeString(result.Value))
		}
		// Join the parts with " | " and add a newline at the end
		line := strings.Join(parts, " | ") + "\n"
		builder.WriteString(line)
	}
	return builder.String()
}

func telegramSearchQuery(ou org_users.OrgUser, sp AiSearchParams) (sql_query_templates.QueryParams, []interface{}) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "telegramSearchQuery"

	args := []interface{}{ou.OrgID}
	q.RawQuery = `SELECT timestamp, group_name, message_text, metadata
				  FROM public.ai_incoming_telegram_msgs
				  WHERE org_id = $1 `

	if sp.SearchContentText != "" {
		args = append(args, sp.SearchContentText)
		q.RawQuery += fmt.Sprintf(` AND message_text_tsvector @@ to_tsquery('english', $%d)`, len(args))
	}
	if sp.GroupFilter != "" {
		args = append(args, sp.GroupFilter)
		q.RawQuery += `AND group_name ILIKE '%' || ` + fmt.Sprintf("$%d", len(args)) + ` || '%' `
	}
	if !sp.SearchInterval[0].IsZero() && !sp.SearchInterval[1].IsZero() {
		if len(args) > 0 {
			q.RawQuery += ` AND`
		} else {
			q.RawQuery += ` WHERE`
		}
		tsRangeStart, tsEnd := sp.SearchInterval.GetUnixTimestamps()
		q.RawQuery += fmt.Sprintf(` timestamp BETWEEN $%d AND $%d `, len(args)+1, len(args)+2)
		args = append(args, tsRangeStart, tsEnd)
	}

	q.RawQuery += ` ORDER BY timestamp DESC;`
	return q, args
}

func SearchTelegram(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q, args := telegramSearchQuery(ou, sp)
	var srs []SearchResult
	rows, err := apps.Pg.Query(ctx, q.RawQuery, args...)
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

const Sn = "OpenAI"

func HashParams(orgID int, hashParams []interface{}) (string, error) {
	hash := sha3.New256()
	for i, v := range hashParams {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		if i == 0 {
			_, _ = hash.Write([]byte(fmt.Sprintf("org-%d", orgID)))
		}
		_, _ = hash.Write(b)
	}
	// Get the resulting encoded byte slice
	sha3v := hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(sha3v)), nil
}

type HashedSearchResult struct {
	SearchAndResultsHash string `json:"searchAndResultHash"`
	SearchAnalysisHash   string `json:"searchAndResultsAndResponseHash"`
}

func HashAiSearchResponseResultsAndParams(ou org_users.OrgUser, response openai.ChatCompletionResponse, sp AiSearchParams, sr []SearchResult) (*HashedSearchResult, error) {
	hash1, err := HashParams(ou.OrgID, []interface{}{sp, sr})
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error hashing search: %s", Sn)
		return nil, err
	}
	hash2, err := HashParams(ou.OrgID, []interface{}{sp, sr, response})
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error hashing search params: %s", Sn)
		return nil, err
	}
	hrp := &HashedSearchResult{
		SearchAndResultsHash: hash1,
		SearchAnalysisHash:   hash2,
	}
	return hrp, nil
}

func insertCompletionResp() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertCompletionResponse"
	q.RawQuery =
		`WITH cte_insert_response AS (
            INSERT INTO completion_responses(org_id, user_id, prompt_tokens, completion_tokens, total_tokens, model, completion_choices)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING response_id
        )
        INSERT INTO ai_search_analysis_results(response_id, search_hash, analysis_hash, search_params, search_results)
        SELECT response_id, $8, $9, $10, $11
        FROM cte_insert_response
        ON CONFLICT (analysis_hash) DO NOTHING
        RETURNING analysis_id;
        `
	return q
}

func SanitizeSearchResults(results []SearchResult) []SearchResult {
	for i, _ := range results {
		results[i].Group = sanitizeUTF8(results[i].Group)
		results[i].Value = sanitizeUTF8(results[i].Value)
		results[i].Metadata.Sanitize()
	}
	return results
}

func SanitizeSearchParams(sp *AiSearchParams) {
	if sp == nil {
		return
	}
	sp.Usernames = sanitizeUTF8(sp.Usernames)
	sp.SearchContentText = sanitizeUTF8(sp.SearchContentText)
	sp.GroupFilter = sanitizeUTF8(sp.GroupFilter)
}

func InsertCompletionResponseChatGptFromSearch(ctx context.Context, ou org_users.OrgUser, response openai.ChatCompletionResponse, sp AiSearchParams, sr []SearchResult) error {
	q := insertCompletionResp()
	for i, choice := range response.Choices {
		response.Choices[i].Message.Content = sanitizeUTF8(choice.Message.Content)
	}
	sr = SanitizeSearchResults(sr)
	SanitizeSearchParams(&sp)
	hrp, err := HashAiSearchResponseResultsAndParams(ou, response, sp, sr)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error hashing search params: %s", Sn)
		return err
	}
	searchParams, err := json.Marshal(sp)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	searchResults, err := json.Marshal(sr)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	completionChoices, err := json.Marshal(response.Choices)
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	r, err := apps.Pg.Exec(ctx, q.RawQuery, ou.OrgID, ou.UserID, response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens, response.Model, completionChoices,
		hrp.SearchAndResultsHash, hrp.SearchAnalysisHash, searchParams, searchResults)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Info().Interface("resp", response).Err(err).Msgf("Error inserting completion response: %s", q.LogHeader(Sn))
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("OrgUser: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return err
}
