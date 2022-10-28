package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type SuperParentClass struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	*ChildClassSingleValue
	*ChildClassMultiValue
}
