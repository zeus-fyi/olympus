package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Orchestrations struct {
	OrchestrationStrID string `db:"-" json:"orchestrationStrID"`
	OrchestrationID    int    `db:"orchestration_id" json:"orchestrationID"`
	OrgID              int    `db:"org_id" json:"orgID"`
	Active             bool   `db:"active" json:"active"`
	GroupName          string `db:"group_name" json:"groupName"`
	Type               string `db:"type" json:"type"`
	Instructions       string `db:"instructions" json:"instructions"`
	OrchestrationName  string `db:"orchestration_name" json:"orchestrationName"`
}
type OrchestrationsSlice []Orchestrations

func (o *Orchestrations) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrchestrationID, o.OrgID, o.OrchestrationName, o.Active, o.Instructions, o.GroupName, o.Type}
	}
	return pgValues
}
func (o *Orchestrations) GetTableColumns() (columnValues []string) {
	columnValues = []string{"orchestration_id", "org_id", "orchestration_name", "active", "instructions", "group_name", "type"}
	return columnValues
}
func (o *Orchestrations) GetTableName() (tableName string) {
	tableName = "orchestrations"
	return tableName
}
