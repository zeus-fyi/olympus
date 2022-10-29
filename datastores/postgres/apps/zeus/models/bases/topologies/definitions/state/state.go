package state

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

// State specific deployed topology to user (statuses can be pending, terminated, etc)
type State struct {
	autogen_bases.TopologiesDeployed
}

func NewState() State {
	s := State{autogen_bases.TopologiesDeployed{
		TopologyStatus: "Pending",
		TopologyID:     0,
		OrgID:          0,
		UserID:         0,
	}}

	return s
}
