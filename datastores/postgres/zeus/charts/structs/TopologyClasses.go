package models

type TopologyClasses struct {
	TopologyClassID     int    `db:"topology_class_id"`
	TopologyClassTypeID int    `db:"topology_class_type_id"`
	TopologyClassName   string `db:"topology_class_name"`
}
