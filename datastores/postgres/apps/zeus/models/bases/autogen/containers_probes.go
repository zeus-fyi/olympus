package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersProbes struct {
	ProbeID     int    `db:"probe_id"`
	ContainerID int    `db:"container_id"`
	ProbeType   string `db:"probe_type"`
}
type ContainersProbesSlice []ContainersProbes

func (c *ContainersProbes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ProbeID, c.ContainerID, c.ProbeType}
	}
	return pgValues
}
func (c *ContainersProbes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"probe_id", "container_id", "probe_type"}
	return columnValues
}
func (c *ContainersProbes) GetTableName() (tableName string) {
	tableName = "containers_probes"
	return tableName
}
