package autogen_structs

type TopologiesDeployed struct {
	TopologyID     int    `db:"topology_id"`
	OrgID          int    `db:"org_id"`
	UserID         int    `db:"user_id"`
	TopologyStatus string `db:"topology_status"`
}
