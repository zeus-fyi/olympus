package common

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type SpecWorkload struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	Replicas ChildClassSingleValue
	Selector Selector
}

func NewSpecWorkload() SpecWorkload {
	sw := SpecWorkload{}
	sw.Selector = NewSelector()
	sw.Replicas = NewInitChildClassSingleValue("replicas", "0")
	return sw
}
