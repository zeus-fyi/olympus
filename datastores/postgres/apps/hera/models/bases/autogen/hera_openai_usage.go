package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type HeraOpenaiUsage struct {
	OrgID           int `db:"org_id" json:"orgID"`
	TokensRemaining int `db:"tokens_remaining" json:"tokensRemaining"`
	TokensConsumed  int `db:"tokens_consumed" json:"tokensConsumed"`
}
type HeraOpenaiUsageSlice []HeraOpenaiUsage

func (h *HeraOpenaiUsage) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{h.OrgID, h.TokensRemaining, h.TokensConsumed}
	}
	return pgValues
}
func (h *HeraOpenaiUsage) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "tokens_remaining", "tokens_consumed"}
	return columnValues
}
func (h *HeraOpenaiUsage) GetTableName() (tableName string) {
	tableName = "hera_openai_usage"
	return tableName
}
