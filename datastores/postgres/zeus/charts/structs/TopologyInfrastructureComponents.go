package models

type TopologyInfrastructureComponents struct {
	TopologyID     int `db:"topology_id"`
	ChartPackageID int `db:"chart_package_id"`
}
