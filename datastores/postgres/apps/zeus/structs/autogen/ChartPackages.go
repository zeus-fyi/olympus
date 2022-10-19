package autogen_structs

import (
	"database/sql"
)

type ChartPackages struct {
	ChartPackageID   int            `db:"chart_package_id"`
	ChartName        string         `db:"chart_name"`
	ChartVersion     string         `db:"chart_version"`
	ChartDescription sql.NullString `db:"chart_description"`
}
