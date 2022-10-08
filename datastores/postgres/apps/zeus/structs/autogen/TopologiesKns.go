package autogen_structs

type TopologiesKns struct {
	TopologyID int    `db:"topology_id"`
	Context    string `db:"context"`
	Namespace  string `db:"namespace"`
	Env        string `db:"env"`
}
