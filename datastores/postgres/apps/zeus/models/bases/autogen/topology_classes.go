package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologyClasses struct {
	TopologyClassID     int    `db:"topology_class_id" json:"topology_class_id"`
	TopologyClassTypeID int    `db:"topology_class_type_id" json:"topology_class_type_id"`
	TopologyClassName   string `db:"topology_class_name" json:"topology_class_name"`
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
