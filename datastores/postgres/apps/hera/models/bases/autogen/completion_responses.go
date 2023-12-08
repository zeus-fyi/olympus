package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type CompletionResponses struct {
	ResponseID        int    `db:"response_id" json:"responseID"`
	OrgID             int    `db:"org_id" json:"orgID"`
	UserID            int    `db:"user_id" json:"userID"`
	PromptTokens      int    `db:"prompt_tokens" json:"promptTokens"`
	CompletionTokens  int    `db:"completion_tokens" json:"completionTokens"`
	TotalTokens       int    `db:"total_tokens" json:"totalTokens"`
	Model             string `db:"model" json:"model"`
	CompletionChoices string `db:"completion_choices" json:"completionChoices"`
	Prompt            []byte `db:"prompt" json:"prompt,omitempty"`
}
type CompletionResponsesSlice []CompletionResponses

func (c *CompletionResponses) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ResponseID, c.OrgID, c.UserID, c.PromptTokens, c.CompletionTokens, c.TotalTokens, c.Model, c.CompletionChoices, c.Prompt}
	}
	return pgValues
}
func (c *CompletionResponses) GetTableColumns() (columnValues []string) {
	columnValues = []string{"response_id", "org_id", "user_id", "prompt_tokens", "completion_tokens", "total_tokens", "model", "completion_choices", "prompt"}
	return columnValues
}
func (c *CompletionResponses) GetTableName() (tableName string) {
	tableName = "completion_responses"
	return tableName
}
