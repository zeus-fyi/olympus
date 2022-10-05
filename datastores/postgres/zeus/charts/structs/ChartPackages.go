package models

type ChartPackages struct {
	ChartPackageID int    `db:"chart_package_id"`
	ChartName      string `db:"chart_name"`
	ChartVersion   string `db:"chart_version"`
}
