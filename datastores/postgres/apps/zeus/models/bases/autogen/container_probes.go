package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerProbes struct {
	ProbeID             int    `db:"probe_id" json:"probeID"`
	ProbeKeyValuesJSONb string `db:"probe_key_values_jsonb" json:"probeKeyValuesJsonb"`
}
type ContainerProbesSlice []ContainerProbes

func (c *ContainerProbes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ProbeID, c.ProbeKeyValuesJSONb}
	}
	return pgValues
}
func (c *ContainerProbes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"probe_id", "probe_key_values_jsonb"}
	return columnValues
}
func (c *ContainerProbes) GetTableName() (tableName string) {
	tableName = "container_probes"
	return tableName
}
