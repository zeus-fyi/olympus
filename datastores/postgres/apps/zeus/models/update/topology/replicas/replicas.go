package update_replicas

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

const Sn = "UpdateReplicaCountDeployment"

var ts chronos.Chronos

type ReplicaUpdate struct {
	org_users.OrgUser
	TopologyID   int
	ReplicaCount string
}
