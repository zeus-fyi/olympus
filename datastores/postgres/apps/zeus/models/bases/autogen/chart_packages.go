package autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type ChartPackages struct {
	ChartPackageID   int            `db:"chart_package_id" json:"chartPackageID"`
	ChartName        string         `db:"chart_name" json:"chartName"`
	ChartVersion     string         `db:"chart_version" json:"chartVersion"`
	ChartDescription sql.NullString `db:"chart_description" json:"chartDescription"`
}
type ChartPackagesSlice []ChartPackages

func (c *ChartPackages) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartPackageID, c.ChartName, c.ChartVersion, c.ChartDescription}
	}
	return pgValues
}
func (c *ChartPackages) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_package_id", "chart_name", "chart_version", "chart_description"}
	return columnValues
}
func (c *ChartPackages) GetTableName() (tableName string) {
	tableName = "chart_packages"
	return tableName
}
