package autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type TopologyClassTypes struct {
	TopologyClassTypeID   int            `db:"topology_class_type_id" json:"topologyClassTypeID"`
	TopologyClassTypeName sql.NullString `db:"topology_class_type_name" json:"topologyClassTypeName"`
}
type TopologyClassTypesSlice []TopologyClassTypes

func (t *TopologyClassTypes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyClassTypeID, t.TopologyClassTypeName}
	}
	return pgValues
}
func (t *TopologyClassTypes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_class_type_id", "topology_class_type_name"}
	return columnValues
}
func (t *TopologyClassTypes) GetTableName() (tableName string) {
	tableName = "topology_class_types"
	return tableName
}
