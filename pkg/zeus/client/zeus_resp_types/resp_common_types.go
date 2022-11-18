package zeus_resp_types

import (
	"database/sql"
)

type TopologyCreateResponse struct {
	TopologyID int `json:"id"`
}

type ReadTopologiesMetadata struct {
	TopologyID       int            `db:"topology_id" json:"topologyID"`
	TopologyName     string         `db:"topology_name" json:"topologyName"`
	ChartName        string         `db:"chart_name" json:"chartName"`
	ChartVersion     string         `db:"chart_version" json:"chartVersion"`
	ChartDescription sql.NullString `db:"chart_description" json:"chartDescription"`
}

type ReadTopologiesMetadataGroup struct {
	Slice []ReadTopologiesMetadata
}
