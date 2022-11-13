package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyClasses struct {
	TopologyClassID     int    `db:"topology_class_id" json:"topologyClassID"`
	TopologyClassTypeID int    `db:"topology_class_type_id" json:"topologyClassTypeID"`
	TopologyClassName   string `db:"topology_class_name" json:"topologyClassName"`
}
type TopologyClassesSlice []TopologyClasses

func (t *TopologyClasses) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyClassID, t.TopologyClassTypeID, t.TopologyClassName}
	}
	return pgValues
}
func (t *TopologyClasses) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_class_id", "topology_class_type_id", "topology_class_name"}
	return columnValues
}
func (t *TopologyClasses) GetTableName() (tableName string) {
	tableName = "topology_classes"
	return tableName
}
