package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologiesDeployed struct {
	TopologyID     int    `db:"topology_id"`
	OrgID          int    `db:"org_id"`
	UserID         int    `db:"user_id"`
	TopologyStatus string `db:"topology_status"`
}
type TopologiesDeployedSlice []TopologiesDeployed

func (t *TopologiesDeployed) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.OrgID, t.UserID, t.TopologyStatus}
	}
	return pgValues
}
func (t *TopologiesDeployed) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "org_id", "user_id", "topology_status"}
	return columnValues
}
func (t *TopologiesDeployed) GetTableName() (tableName string) {
	tableName = "topologies_deployed"
	return tableName
}
