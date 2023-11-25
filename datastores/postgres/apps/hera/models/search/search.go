package hera_search

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

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
	SearchContentText    string `json:"searchContentText,omitempty"`
	GroupFilter          string `json:"groupFilter,omitempty"`
	Platforms            string `json:"platforms,omitempty"`
	Usernames            string `json:"usernames,omitempty"`
	WorkflowInstructions string `json:"workflowInstructions,omitempty"`
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

type DiscordMetadata struct {
	GuildName    string `json:"guildName"`
	Category     string `json:"topic"`
	CategoryName string `json:"categoryName"`
}

func discordSearchQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "discordSearchQuery"
	q.RawQuery = `SELECT cm.timestamp_creation, cm.content, gi.name, ci.category, ci.name
				  FROM public.ai_incoming_discord_messages cm
				  JOIN public.ai_discord_channel ci ON ci.channel_id = cm.channel_id
				  JOIN public.ai_discord_guild gi ON gi.guild_id = cm.guild_id
        		  WHERE content_tsvector @@ to_tsquery('english', $1)
				  ORDER BY cm.timestamp_creation DESC;`
	return q
}

func discordSearchQuery2() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "discordSearchQuery2"
	q.RawQuery = `SELECT cm.timestamp_creation, cm.content, gi.name, ci.category, ci.name
				  FROM public.ai_incoming_discord_messages cm
				  JOIN public.ai_discord_channel ci ON ci.channel_id = cm.channel_id
				  JOIN public.ai_discord_guild gi ON gi.guild_id = cm.guild_id
				  ORDER BY cm.timestamp_creation DESC;`
	return q
}

//func discordSearchQuery3() sql_query_templates.QueryParams {
//	q := sql_query_templates.QueryParams{}
//	q.QueryName = "discordSearchQuery3"
//	q.RawQuery = `SELECT cm.timestamp_creation, cm.content, gi.name, ci.category, ci.name
//				  FROM public.ai_incoming_discord_messages cm
//				  JOIN public.ai_discord_channel ci ON ci.channel_id = cm.channel_id
//				  JOIN public.ai_discord_guild gi ON gi.guild_id = cm.guild_id
//				  ORDER BY cm.timestamp_creation DESC;`
//	return q
//}

func SearchDiscord(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q := discordSearchQuery()
	var srs []SearchResult

	var rows pgx.Rows
	var err error
	if sp.SearchContentText == "" {
		q = discordSearchQuery2()
		rows, err = apps.Pg.Query(ctx, q.RawQuery)
	} else {
		rows, err = apps.Pg.Query(ctx, q.RawQuery, sp.SearchContentText)
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SearchDiscord")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sr := SearchResult{Source: "discord"}
		rowErr := rows.Scan(&sr.UnixTimestamp, &sr.Value, &sr.Group, &sr.DiscordMetadata.Category, &sr.DiscordMetadata.CategoryName)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SearchDiscord"))
			return nil, rowErr
		}
		srs = append(srs, sr)
	}
	return srs, nil
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

func telegramSearchQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "telegramSearchQuery"
	q.RawQuery = `SELECT timestamp, group_name, message_text, metadata
				  FROM public.ai_incoming_telegram_msgs
              	  WHERE org_id = $1 AND group_name ILIKE '%' || $2 || '%' 
				  ORDER BY chat_id, message_id DESC;`
	return q
}

func telegramSearchQueryWithContent() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "telegramSearchQueryWithContent"
	q.RawQuery = `SELECT timestamp, group_name, message_text, metadata
				  FROM public.ai_incoming_telegram_msgs
              	  WHERE org_id = $1 AND group_name ILIKE '%' || $2 || '%' 
              	  AND message_text_tsvector @@ to_tsquery('english', $3)
				  ORDER BY chat_id, message_id DESC;`
	return q
}

func SearchTelegram(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	q := telegramSearchQuery()
	var srs []SearchResult
	args := []interface{}{ou.OrgID, sp.GroupFilter}
	if sp.SearchContentText != "" {
		args = append(args, sp.SearchContentText)
		q = telegramSearchQueryWithContent()
	}
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

/*
CREATE TABLE public.ai_search_analysis_results(
    analysis_id int8 NOT NULL DEFAULT next_id(),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id),
    search_hash text NOT NULL, -- hash of the search parameters and results
    analysis_hash text NOT NULL, -- hash of the search parameters, results, and response
    search_params JSONB NOT NULL,
    search_results JSONB NOT NULL,
    UNIQUE(analysis_hash)
);
*/

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
	sp.WorkflowInstructions = sanitizeUTF8(sp.WorkflowInstructions)
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
