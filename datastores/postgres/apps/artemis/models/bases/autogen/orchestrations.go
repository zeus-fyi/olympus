package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Orchestrations struct {
	OrchestrationID   int    `db:"orchestration_id" json:"orchestrationID"`
	OrgID             int    `db:"org_id" json:"orgID"`
	Active            bool   `db:"active" json:"active"`
	Instructions      string `db:"instructions" json:"instructions"`
	OrchestrationName string `db:"orchestration_name" json:"orchestrationName"`
}
type OrchestrationsSlice []Orchestrations

func (o *Orchestrations) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrchestrationID, o.OrgID, o.OrchestrationName, o.Active, o.Instructions}
	}
	return pgValues
}
func (o *Orchestrations) GetTableColumns() (columnValues []string) {
	columnValues = []string{"orchestration_id", "org_id", "orchestration_name", "active", "instructions"}
	return columnValues
}
func (o *Orchestrations) GetTableName() (tableName string) {
	tableName = "orchestrations"
	return tableName
}
