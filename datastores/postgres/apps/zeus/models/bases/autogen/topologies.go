package autogen_bases

type Topologies struct {
	TopologyID int    `db:"topology_id"`
	Name       string `db:"name"`
}
type TopologiesSlice []Topologies

func (t *Topologies) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologyID, t.Name}
	}
	return pgValues
}
func (t *Topologies) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_id", "name"}
	return columnValues
}
func (t *Topologies) GetTableName() (tableName string) {
	tableName = "topologies"
	return tableName
}
