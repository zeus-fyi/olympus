package hera_search

import "context"

type AiSearchParams struct {
	SearchContentText    string `json:"searchContentText,omitempty"`
	GroupFilter          string `json:"groupFilter,omitempty"`
	Usernames            string `json:"usernames,omitempty"`
	WorkflowInstructions string `json:"workflowInstructions,omitempty"`
}

type SearchResult struct {
	Value string `json:"value,omitempty"`
}

func SearchTelegram(ctx context.Context, sp AiSearchParams) ([]SearchResult, error) {

	return nil, nil
}
