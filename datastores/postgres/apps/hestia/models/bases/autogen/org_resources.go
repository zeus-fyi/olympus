package hestia_autogen_bases

import (
	"database/sql"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type OrgResources struct {
	EndService   sql.NullTime `db:"end_service" json:"endService"`
	ResourceID   int          `db:"resource_id" json:"rID"`
	OrgID        int          `db:"org_id" json:"orgID"`
	BeginService time.Time    `db:"begin_service" json:"beginService"`
	Quantity     float64      `db:"quantity" json:"quantity"`
	FreeTrial    bool         `db:"free_trial" json:"freeTrial"`
}
type OrgResourcesSlice []OrgResources

func (o *OrgResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.EndService, o.ResourceID, o.OrgID, o.BeginService, o.Quantity, o.FreeTrial}
	}
	return pgValues
}
func (o *OrgResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"end_service", "resource_id", "org_id", "begin_service", "quantity", "free_trial"}
	return columnValues
}
func (o *OrgResources) GetTableName() (tableName string) {
	tableName = "org_resources"
	return tableName
}
