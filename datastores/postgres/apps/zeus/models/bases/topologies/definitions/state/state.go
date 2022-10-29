package state

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type State struct {
	autogen_bases.TopologiesDeployed
}

func NewState() State {
	s := State{autogen_bases.TopologiesDeployed{
		TopologyStatus: "",
		TopologyID:     0,
		OrgID:          0,
		UserID:         0,
	}}

	return s
}
