package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type TopologySystemComponentsSecrets struct {
	TopologySystemComponentID int `db:"topology_system_component_id" json:"topologySystemComponentID"`
	SecretID                  int `db:"secret_id" json:"secretID"`
}
type TopologySystemComponentsSecretsSlice []TopologySystemComponentsSecrets

func (t *TopologySystemComponentsSecrets) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.TopologySystemComponentID, t.SecretID}
	}
	return pgValues
}
func (t *TopologySystemComponentsSecrets) GetTableColumns() (columnValues []string) {
	columnValues = []string{"topology_system_component_id", "secret_id"}
	return columnValues
}
func (t *TopologySystemComponentsSecrets) GetTableName() (tableName string) {
	tableName = "topology_system_components_secrets"
	return tableName
}
