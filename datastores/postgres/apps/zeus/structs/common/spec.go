package common

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type SpecWorkload struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	Replicas ChildClassSingleValue
	Selector Selector
}
