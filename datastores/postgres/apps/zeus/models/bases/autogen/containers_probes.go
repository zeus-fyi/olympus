package autogen_bases

type ContainersProbes struct {
	ProbeType   string `db:"probe_type" json:"probeType"`
	ProbeID     int    `db:"probe_id" json:"probeID"`
	ContainerID int    `db:"container_id" json:"containerID"`
}
type ContainersProbesSlice []ContainersProbes

func (c *ContainersProbes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ProbeType, c.ProbeID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersProbes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"probe_type", "probe_id", "container_id"}
	return columnValues
}
func (c *ContainersProbes) GetTableName() (tableName string) {
	tableName = "containers_probes"
	return tableName
}
