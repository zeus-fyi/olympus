package autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type TopologiesDeployed struct {
	TopologyID     int       `db:"topology_id" json:"topology_id"`
	TopologyStatus string    `db:"topology_status" json:"topology_status"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
type TopologiesDeployedSlice []TopologiesDeployed

func (t *TopologiesDeployed) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.TopologyStatus, t.UpdatedAt}
	}
	return pgValues
}
func (t *TopologiesDeployed) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "topology_status", "updated_at"}
	return columnValues
}
func (t *TopologiesDeployed) GetTableName() (tableName string) {
	tableName = "topologies_deployed"
	return tableName
}
