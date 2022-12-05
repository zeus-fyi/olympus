package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyBaseComponents struct {
	TopologyBaseName          string `db:"topology_base_name" json:"topologyBaseName"`
	OrgID                     int    `db:"org_id" json:"orgID"`
	TopologySystemComponentID int    `db:"topology_system_component_id" json:"topologySystemComponentID"`
	TopologyClassTypeID       int    `db:"topology_class_type_id" json:"topologyClassTypeID"`
	TopologyBaseComponentID   int    `db:"topology_base_component_id" json:"topologyBaseComponentID"`
}
type TopologyBaseComponentsSlice []TopologyBaseComponents

func (t *TopologyBaseComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyBaseName, t.OrgID, t.TopologySystemComponentID, t.TopologyClassTypeID, t.TopologyBaseComponentID}
	}
	return pgValues
}
func (t *TopologyBaseComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_base_name", "org_id", "topology_system_component_id", "topology_class_type_id", "topology_base_component_id"}
	return columnValues
}
func (t *TopologyBaseComponents) GetTableName() (tableName string) {
	tableName = "topology_base_components"
	return tableName
}
