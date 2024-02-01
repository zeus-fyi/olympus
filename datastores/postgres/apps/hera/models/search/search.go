package hera_search

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"golang.org/x/crypto/sha3"
)

type AiSearchParams struct {
	Retrieval artemis_orchestrations.RetrievalItem `json:"retrieval,omitempty"`
	TimeRange string                               `json:"timeRange,omitempty"`
	Window    artemis_orchestrations.Window        `json:"window,omitempty"`
}

func TimeRangeStringToWindow(sp *AiSearchParams) {
	if sp == nil {
		return
	}
	ts := time.Now()
	w := artemis_orchestrations.Window{}
	switch sp.TimeRange {
	case "1 hour":
		w.Start = ts.Add(-1 * time.Hour)
		w.End = ts
	case "24 hours":
		w.Start = ts.AddDate(0, 0, -1)
		w.End = ts
	case "7 days":
		w.Start = ts.AddDate(0, 0, -7)
		w.End = ts
	case "30 days":
		w.Start = ts.AddDate(0, 0, -30)
		w.End = ts
	case "all":
		w.Start = time.Unix(0, 0)
		w.End = ts
	case "window":
		log.Info().Interface("searchInterval", w).Msg("window")
	}

	w.UnixStartTime = int(w.Start.Unix())
	w.UnixEndTime = int(w.End.Unix())
	sp.Window = w
}

type AiModelParams struct {
	Model         string `json:"model"`
	TokenCountMax int    `json:"tokenCountMax"`
}

type SearchResultGroup struct {
	DataIn                         []artemis_orchestrations.AIWorkflowAnalysisResult `json:"dataIn,omitempty"`
	PlatformName                   string                                            `json:"platformName"`
	SourceTaskID                   int                                               `json:"sourceTaskID,omitempty"`
	ExtractionPromptExt            string                                            `json:"extractionPromptExt,omitempty"`
	Model                          string                                            `json:"model,omitempty"`
	ResponseFormat                 string                                            `json:"responseFormat,omitempty"`
	BodyPrompt                     string                                            `json:"bodyPrompt,omitempty"`
	ApiResponseResults             []SearchResult                                    `json:"apiResponseResults,omitempty"`
	SearchResults                  []SearchResult                                    `json:"searchResults"`
	FilteredSearchResults          []SearchResult                                    `json:"filteredSearchResults,omitempty"`
	FilteredSearchResultMap        map[int]*SearchResult                             `json:"filteredSearchResultsMap"`
	SearchResultChunkTokenEstimate *int                                              `json:"searchResultChunkTokenEstimates,omitempty"`
	Window                         artemis_orchestrations.Window                     `json:"window,omitempty"`
	FunctionDefinition             openai.FunctionDefinition                         `json:"functionDefinition,omitempty"`
}

func (sg *SearchResultGroup) GetMessageMap() map[int]*SearchResult {
	msgMap := make(map[int]*SearchResult)
	for _, v := range sg.SearchResults {
		msgMap[v.UnixTimestamp] = &v
	}
	sg.FilteredSearchResultMap = make(map[int]*SearchResult)
	return msgMap
}

func (sg *SearchResultGroup) GetPromptBody() string {
	if len(sg.SearchResults) == 0 || len(sg.ApiResponseResults) == 0 {
		return sg.BodyPrompt
	}
	if len(sg.ApiResponseResults) > 0 {
		return SearchResultSliceToString(sg.ApiResponseResults)
	}
	return FormatSearchResultsV4(sg.FilteredSearchResultMap, sg.SearchResults)
}

func SearchResultSliceToString(results []SearchResult) string {
	var sb strings.Builder

	for _, result := range results {
		if result.WebResponse.Body == nil {
			continue
		}
		bodyString, err := json.Marshal(result.WebResponse.Body)
		if err != nil {
			// Handle error, maybe log it or use a default error message in place of the body
			continue // or handle it differently
		}

		sb.WriteString(string(bodyString))
		sb.WriteString("\n") // Add a newline after each result's body
	}

	return sb.String()
}

func FormatSearchResultsV4(filteredMap map[int]*SearchResult, results []SearchResult) string {
	if len(results) == 0 {
		return ""
	}
	var newResults []SimplifiedSearchResultJSON
	for i, r := range results {
		if fr, ok := filteredMap[r.UnixTimestamp]; ok {
			if fr != nil && fr.Verified != nil && *fr.Verified {
				results[i] = *fr
			}
		}
		if results[i].Verified != nil && *results[i].Verified {
			continue
		}
		nr := SimplifiedSearchResultJSON{
			MessageID:   fmt.Sprintf("%d", r.UnixTimestamp),
			MessageBody: r.Value,
		}
		newResults = append(newResults, nr)
	}
	b, err := json.Marshal(newResults)
	if err != nil {
		log.Err(err).Msg("FormatSearchResultsV3: Error marshalling search results")
		return ""
	}
	return string(b)
}

type SearchResult struct {
	UnixTimestamp   int              `json:"unixTimestamp"`
	Source          string           `json:"source"`
	Value           string           `json:"value"`
	Group           string           `json:"group"`
	Verified        *bool            `json:"verified,omitempty"`
	Metadata        TelegramMetadata `json:"metadata,omitempty"`
	DiscordMetadata DiscordMetadata  `json:"discordMetadata"`
	WebResponse     WebResponse      `json:"webResponses,omitempty"`
}

type WebResponse struct {
	WebFilters *artemis_orchestrations.WebFilters `json:"webFilters,omitempty"`
	Body       echo.Map                           `json:"body"`
	RawMessage []byte                             `json:"rawMessage"`
}

type ByTimestamp []SearchResult

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].UnixTimestamp > a[j].UnixTimestamp }

// SortSearchResults sorts the slice of SearchResult in descending order by UnixTimestamp.
func SortSearchResults(results []SearchResult) {
	sort.Sort(ByTimestamp(results))
}

type SearchIndexerParams struct {
	OrgID           int    `json:"orgID,omitempty"`
	SearchID        int    `json:"searchID"`
	SearchGroupName string `json:"searchGroupName"`
	MaxResults      int    `json:"maxResults"`
	Query           string `json:"query"`
	Platform        string `json:"platform"`
	Active          bool   `json:"active"`
}

func GetSearchIndexersByOrg(ctx context.Context, ou org_users.OrgUser) ([]SearchIndexerParams, error) {
	query := `
		SELECT search_id, search_group_name, max_results, query, 'reddit' AS platform, active
		FROM public.ai_reddit_search_query
		WHERE org_id = $1
		UNION
		SELECT search_id, search_group_name, max_results, query, 'twitter' AS platform, active
		FROM public.ai_twitter_search_query
		WHERE org_id = $1
		UNION
		SELECT
		    dsq.search_id,
		    dsq.search_group_name,
		    dsq.max_results,
		    gi.name || ' | ' || ci.category || ' | ' || ci.name || ' | ' || ci.channel_id AS query,
		    'discord' AS platform, dsq.active
		FROM
		    (
		        SELECT dm.search_id, dm.guild_id, dm.channel_id, dsq.active, MAX(dm.timestamp_creation) AS max_message_id
		        FROM public.ai_incoming_discord_messages dm
		        INNER JOIN public.ai_discord_search_query dsq
		        ON dm.search_id = dsq.search_id
		        WHERE dsq.org_id = $1
		        GROUP BY dm.search_id, dm.guild_id, dm.channel_id, dsq.active
		    ) AS latest_discord_messages
		JOIN
		    public.ai_discord_search_query dsq ON dsq.search_id = latest_discord_messages.search_id
		JOIN
		    public.ai_discord_channel ci ON ci.channel_id = latest_discord_messages.channel_id
		JOIN
		    public.ai_discord_guild gi ON gi.guild_id = latest_discord_messages.guild_id
	`
	rows, err := apps.Pg.Query(ctx, query, ou.OrgID)
	if err != nil {
		log.Err(err).Msg("Error querying search indexers")
		return nil, err
	}
	defer rows.Close()

	var srs []SearchIndexerParams
	for rows.Next() {
		var si SearchIndexerParams
		err = rows.Scan(&si.SearchID, &si.SearchGroupName, &si.MaxResults, &si.Query, &si.Platform, &si.Active)
		if err != nil {
			log.Err(err).Msg("Error querying search indexers")
			return nil, err
		}
		srs = append(srs, si)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return srs, nil
}

func GetAllActiveSearchIndexers(ctx context.Context) ([]SearchIndexerParams, error) {
	query := `
		SELECT search_id, search_group_name, max_results, query, 'reddit' AS platform, active, org_id
		FROM public.ai_reddit_search_query
		WHERE active = true
		UNION
		SELECT search_id, search_group_name, max_results, query, 'twitter' AS platform, active, org_id
		FROM public.ai_twitter_search_query
		WHERE active = true
		UNION
		SELECT
		    dsq.search_id,
		    dsq.search_group_name,
		    dsq.max_results,
		    gi.name || ' | ' || ci.category || ' | ' || ci.name || ' | ' || ci.channel_id AS query,
		    'discord' AS platform,
		    dsq.active,
		    dsq.org_id
		FROM
		    (
		        SELECT dm.search_id, dm.guild_id, dm.channel_id, dsq.active, MAX(dm.timestamp_creation) AS max_message_id
		        FROM public.ai_incoming_discord_messages dm
		        INNER JOIN public.ai_discord_search_query dsq
		        ON dm.search_id = dsq.search_id
		        WHERE dsq.active = true
		        GROUP BY dm.search_id, dm.guild_id, dm.channel_id, dsq.active, dsq.org_id
		    ) AS latest_discord_messages
		JOIN
		    public.ai_discord_search_query dsq ON dsq.search_id = latest_discord_messages.search_id
		JOIN
		    public.ai_discord_channel ci ON ci.channel_id = latest_discord_messages.channel_id
		JOIN
		    public.ai_discord_guild gi ON gi.guild_id = latest_discord_messages.guild_id
	`
	rows, err := apps.Pg.Query(ctx, query)
	if err != nil {
		log.Err(err).Msg("Error querying search indexers")
		return nil, err
	}
	defer rows.Close()
	var srs []SearchIndexerParams
	for rows.Next() {
		var si SearchIndexerParams
		err = rows.Scan(&si.SearchID, &si.SearchGroupName, &si.MaxResults, &si.Query, &si.Platform, &si.Active, &si.OrgID)
		if err != nil {
			log.Err(err).Msg("Error querying search indexers")
			return nil, err
		}
		srs = append(srs, si)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return srs, nil
}

func PerformPlatformSearches(ctx context.Context, ou org_users.OrgUser, sp AiSearchParams) ([]SearchResult, error) {
	var res []SearchResult
	platform := sp.Retrieval.RetrievalPlatform
	if strings.Contains(platform, "twitter") || len(platform) == 0 {
		resTwitter, err := SearchTwitter(ctx, ou, sp)
		if err != nil {
			return nil, err
		}
		res = append(res, resTwitter...)
	}

	if strings.Contains(platform, "discord") || len(platform) == 0 {
		resDiscord, err := SearchDiscord(ctx, ou, sp)
		if err != nil {
			return nil, err
		}
		res = append(res, resDiscord...)
	}

	if strings.Contains(platform, "telegram") || len(platform) == 0 {
		resTelegram, err := SearchTelegram(ctx, ou, sp)
		if err != nil {
			return nil, err
		}
		res = append(res, resTelegram...)
	}

	if strings.Contains(platform, "reddit") || len(platform) == 0 {
		resReddit, err := SearchReddit(ctx, ou, sp)
		if err != nil {
			return nil, err
		}
		res = append(res, resReddit...)
	}

	if strings.Contains(platform, "web") {
		resWeb, err := SearchReddit(ctx, ou, sp)
		if err != nil {
			return nil, err
		}
		res = append(res, resWeb...)
	}
	return res, nil
}

type SearchResults struct {
	AiModelParams AiModelParams  `json:"aiModelParams"`
	Results       []SearchResult `json:"results"`
}

func FormatSearchResultsV2(results []SearchResult) string {
	if len(results) == 0 {
		return ""
	}
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

type SimplifiedSearchResultJSON struct {
	MessageID   string `json:"msg_id"`
	MessageBody string `json:"msg_body"`
}

func FormatSearchResultsV3(results []SearchResult) string {
	if len(results) == 0 {
		return ""
	}
	var newResults []SimplifiedSearchResultJSON
	for _, r := range results {
		if r.Verified != nil && *r.Verified {
			continue
		}
		nr := SimplifiedSearchResultJSON{
			MessageID:   fmt.Sprintf("%d", r.UnixTimestamp),
			MessageBody: r.Value,
		}
		newResults = append(newResults, nr)
	}
	b, err := json.Marshal(newResults)
	if err != nil {
		log.Err(err).Msg("FormatSearchResultsV3: Error marshalling search results")
		return ""
	}
	return string(b)
}

func telegramSearchQuery(ou org_users.OrgUser, sp AiSearchParams) (sql_query_templates.QueryParams, []interface{}) {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "telegramSearchQuery"

	args := []interface{}{ou.OrgID}
	q.RawQuery = `SELECT timestamp, group_name, message_text, metadata
				  FROM public.ai_incoming_telegram_msgs
				  WHERE org_id = $1 `

	if sp.Retrieval.RetrievalKeywords != nil && *sp.Retrieval.RetrievalKeywords != "" {
		args = append(args, sp.Retrieval.RetrievalKeywords)
		q.RawQuery += fmt.Sprintf(` AND message_text_tsvector @@ to_tsquery('english', $%d)`, len(args))
	}
	if sp.Retrieval.RetrievalGroup != "" {
		args = append(args, sp.Retrieval.RetrievalGroup)
		q.RawQuery += `AND group_name ILIKE '%' || ` + fmt.Sprintf("$%d", len(args)) + ` || '%' `
	}
	if !sp.Window.Start.IsZero() && !sp.Window.End.IsZero() {
		if len(args) > 0 {
			q.RawQuery += ` AND`
		} else {
			q.RawQuery += ` WHERE`
		}
		tsRangeStart, tsEnd := sp.Window.GetUnixTimestamps()
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

	if sp.Retrieval.RetrievalUsernames != nil {
		sp.Retrieval.RetrievalUsernames = new(string)
		tmp := sanitizeUTF8(*sp.Retrieval.RetrievalUsernames)
		sp.Retrieval.RetrievalUsernames = &tmp
	}
	if sp.Retrieval.RetrievalKeywords != nil {
		sp.Retrieval.RetrievalKeywords = new(string)
		tmp := sanitizeUTF8(*sp.Retrieval.RetrievalKeywords)
		sp.Retrieval.RetrievalKeywords = &tmp
	}
	if sp.Retrieval.RetrievalPrompt != nil {
		sp.Retrieval.RetrievalPrompt = new(string)
		tmp := sanitizeUTF8(*sp.Retrieval.RetrievalPrompt)
		sp.Retrieval.RetrievalPrompt = &tmp
	}
	if sp.Retrieval.RetrievalPlatformGroups != nil {
		sp.Retrieval.RetrievalPlatformGroups = new(string)
		tmp := sanitizeUTF8(*sp.Retrieval.RetrievalPlatformGroups)
		sp.Retrieval.RetrievalPlatformGroups = &tmp
	}
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
