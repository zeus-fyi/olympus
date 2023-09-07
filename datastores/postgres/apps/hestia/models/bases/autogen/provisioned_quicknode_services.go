package hestia_autogen_bases

import (
	"database/sql"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type ProvisionedQuickNodeServices struct {
	CreatedAt   time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updatedAt"`
	QuickNodeID string         `db:"quicknode_id" json:"quicknodeID"`
	EndpointID  string         `db:"endpoint_id" json:"endpointID"`
	HttpURL     sql.NullString `db:"http_url" json:"httpUrl"`
	Network     sql.NullString `db:"network" json:"network"`
	Plan        string         `db:"plan" json:"plan"`
	Active      bool           `db:"active" json:"active"`
	OrgID       int            `db:"org_id" json:"orgID"`
	WssURL      sql.NullString `db:"wss_url" json:"wssUrl"`
	Chain       sql.NullString `db:"chain" json:"chain"`
}
type ProvisionedQuickNodeServicesSlice []ProvisionedQuickNodeServices

func (p *ProvisionedQuickNodeServices) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{p.CreatedAt, p.UpdatedAt, p.QuickNodeID, p.EndpointID, p.HttpURL, p.Network, p.Plan, p.Active, p.OrgID, p.WssURL, p.Chain}
	}
	return pgValues
}
func (p *ProvisionedQuickNodeServices) GetTableColumns() (columnValues []string) {
	columnValues = []string{"created_at", "updated_at", "quicknode_id", "endpoint_id", "http_url", "network", "plan", "active", "org_id", "wss_url", "chain"}
	return columnValues
}
func (p *ProvisionedQuickNodeServices) GetTableName() (tableName string) {
	tableName = "provisioned_quicknode_services"
	return tableName
}
