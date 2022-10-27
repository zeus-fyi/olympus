package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologiesDeployed struct {
	TopologyStatus string `db:"topology_status" json:"topology_status"`
	TopologyID     int    `db:"topology_id" json:"topology_id"`
	OrgID          int    `db:"org_id" json:"org_id"`
	UserID         int    `db:"user_id" json:"user_id"`
}
type TopologiesDeployedSlice []TopologiesDeployed

func (t *TopologiesDeployed) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyStatus, t.TopologyID, t.OrgID, t.UserID}
	}
	return pgValues
}
func (t *TopologiesDeployed) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_status", "topology_id", "org_id", "user_id"}
	return columnValues
}
func (t *TopologiesDeployed) GetTableName() (tableName string) {
	tableName = "topologies_deployed"
	return tableName
}
