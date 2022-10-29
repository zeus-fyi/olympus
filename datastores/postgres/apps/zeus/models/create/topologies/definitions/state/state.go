package create_state

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"

type State struct {
	state.State
}

func NewCreateState() State {
	s := state.NewState()
	cs := State{s}
	return cs
}
