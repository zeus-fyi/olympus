package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologySystemComponents struct {
	OrgID                       int    `db:"org_id" json:"orgID"`
	TopologySystemComponentID   int    `db:"topology_system_component_id" json:"topologySystemComponentID"`
	TopologyClassTypeID         int    `db:"topology_class_type_id" json:"topologyClassTypeID"`
	TopologySystemComponentName string `db:"topology_system_component_name" json:"topologySystemComponentName"`
}
type TopologySystemComponentsSlice []TopologySystemComponents

func (t *TopologySystemComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.OrgID, t.TopologySystemComponentID, t.TopologyClassTypeID, t.TopologySystemComponentName}
	}
	return pgValues
}
func (t *TopologySystemComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "topology_system_component_id", "topology_class_type_id", "topology_system_component_name"}
	return columnValues
}
func (t *TopologySystemComponents) GetTableName() (tableName string) {
	tableName = "topology_system_components"
	return tableName
}
