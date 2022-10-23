package autogen_bases

type TopologiesDeployed struct {
	OrgID          int    `db:"org_id"`
	UserID         int    `db:"user_id"`
	TopologyStatus string `db:"topology_status"`
	TopologyID     int    `db:"topology_id"`
}
type TopologiesDeployedSlice []TopologiesDeployed

func (t *TopologiesDeployed) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.OrgID, t.UserID, t.TopologyStatus, t.TopologyID}
	}
	return pgValues
}
func (t *TopologiesDeployed) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "user_id", "topology_status", "topology_id"}
	return columnValues
}
func (t *TopologiesDeployed) GetTableName() (tableName string) {
	tableName = "topologies_deployed"
	return tableName
}
