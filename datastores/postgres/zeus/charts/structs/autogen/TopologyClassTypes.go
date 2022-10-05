package autogen_structs

import (
	"database/sql"
)

type TopologyClassTypes struct {
	TopologyClassTypeID   int            `db:"topology_class_type_id"`
	TopologyClassTypeName sql.NullString `db:"topology_class_type_name"`
}
