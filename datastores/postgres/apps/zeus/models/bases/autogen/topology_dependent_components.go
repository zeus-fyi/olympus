package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyDependentComponents struct {
	TopologyClassID int `db:"topology_class_id"`
	TopologyID      int `db:"topology_id"`
}
type TopologyDependentComponentsSlice []TopologyDependentComponents

func (t *TopologyDependentComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyClassID, t.TopologyID}
	}
	return pgValues
}
func (t *TopologyDependentComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_class_id", "topology_id"}
	return columnValues
}
func (t *TopologyDependentComponents) GetTableName() (tableName string) {
	tableName = "topology_dependent_components"
	return tableName
}
