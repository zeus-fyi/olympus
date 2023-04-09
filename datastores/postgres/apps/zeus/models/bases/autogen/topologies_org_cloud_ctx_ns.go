package autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type TopologiesOrgCloudCtxNs struct {
	CloudCtxNsID   int       `db:"cloud_ctx_ns_id" json:"cloudCtxNsID"`
	OrgID          int       `db:"org_id" json:"orgID,omitempty"`
	CloudProvider  string    `db:"cloud_provider" json:"cloudProvider"`
	Context        string    `db:"context" json:"context"`
	Region         string    `db:"region" json:"region"`
	Namespace      string    `db:"namespace" json:"namespace"`
	NamespaceAlias string    `db:"namespace_alias" json:"namespaceAlias"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}
type TopologiesOrgCloudCtxNsSlice []TopologiesOrgCloudCtxNs

func (t *TopologiesOrgCloudCtxNs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{t.CloudCtxNsID, t.OrgID, t.CloudProvider, t.Context, t.Region, t.Namespace, t.NamespaceAlias}
	}
	return pgValues
}
func (t *TopologiesOrgCloudCtxNs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"cloud_ctx_ns_id", "org_id", "cloud_provider", "context", "region", "namespace", "namespace_alias", "created_at"}
	return columnValues
}
func (t *TopologiesOrgCloudCtxNs) GetTableName() (tableName string) {
	tableName = "topologies_org_cloud_ctx_ns"
	return tableName
}
