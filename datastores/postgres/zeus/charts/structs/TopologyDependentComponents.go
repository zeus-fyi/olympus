package models

type TopologyDependentComponents struct {
	TopologyClassID int `db:"topology_class_id"`
	TopologyID      int `db:"topology_id"`
}
