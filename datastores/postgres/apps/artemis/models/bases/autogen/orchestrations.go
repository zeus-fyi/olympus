package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Orchestrations struct {
	OrchestrationID   int    `db:"orchestration_id" json:"orchestrationID"`
	OrgID             int    `db:"org_id" json:"orgID"`
	OrchestrationName string `db:"orchestration_name" json:"orchestrationName"`
}
type OrchestrationsSlice []Orchestrations

func (o *Orchestrations) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrchestrationID, o.OrgID, o.OrchestrationName}
	}
	return pgValues
}
func (o *Orchestrations) GetTableColumns() (columnValues []string) {
	columnValues = []string{"orchestration_id", "org_id", "orchestration_name"}
	return columnValues
}
func (o *Orchestrations) GetTableName() (tableName string) {
	tableName = "orchestrations"
	return tableName
}
