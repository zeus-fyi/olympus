package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Topologies struct {
	Name       string `db:"name" json:"name"`
	TopologyID int    `db:"topology_id" json:"topology_id"`
}
type TopologiesSlice []Topologies

func (t *Topologies) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.Name, t.TopologyID}
	}
	return pgValues
}
func (t *Topologies) GetTableColumns() (columnValues []string) {
	columnValues = []string{"name", "topology_id"}
	return columnValues
}
func (t *Topologies) GetTableName() (tableName string) {
	tableName = "topologies"
	return tableName
}
