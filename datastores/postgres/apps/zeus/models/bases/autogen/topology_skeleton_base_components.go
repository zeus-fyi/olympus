package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologySkeletonBaseComponents struct {
	OrgID                    int    `db:"org_id" json:"orgID"`
	TopologyBaseComponentID  int    `db:"topology_base_component_id" json:"topologyBaseComponentID"`
	TopologyClassTypeID      int    `db:"topology_class_type_id" json:"topologyClassTypeID"`
	TopologySkeletonBaseID   int    `db:"topology_skeleton_base_id" json:"topologySkeletonBaseID"`
	TopologySkeletonBaseName string `db:"topology_skeleton_base_name" json:"topologySkeletonBaseName"`
}
type TopologySkeletonBaseComponentsSlice []TopologySkeletonBaseComponents

func (t *TopologySkeletonBaseComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.OrgID, t.TopologyBaseComponentID, t.TopologyClassTypeID, t.TopologySkeletonBaseID, t.TopologySkeletonBaseName}
	}
	return pgValues
}
func (t *TopologySkeletonBaseComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "topology_base_component_id", "topology_class_type_id", "topology_skeleton_base_id", "topology_skeleton_base_name"}
	return columnValues
}
func (t *TopologySkeletonBaseComponents) GetTableName() (tableName string) {
	tableName = "topology_skeleton_base_components"
	return tableName
}
