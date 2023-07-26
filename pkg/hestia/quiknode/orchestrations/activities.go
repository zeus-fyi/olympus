package quicknode_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/quiknode"
)

type HestiaQuicknodeActivities struct {
}

func NewHestiaQuicknodeActivities() HestiaQuicknodeActivities {
	return HestiaQuicknodeActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaQuicknodeActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		h.Provision, h.UpdateProvision, h.Deprovision, h.Deprovision,
	}
}

func (h *HestiaQuicknodeActivities) Provision(pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {

	return nil
}

func (h *HestiaQuicknodeActivities) UpdateProvision(pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {

	return nil
}

func (h *HestiaQuicknodeActivities) Deprovision(dp hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {

	return nil
}

func (h *HestiaQuicknodeActivities) Deactivate(da hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {

	return nil
}
