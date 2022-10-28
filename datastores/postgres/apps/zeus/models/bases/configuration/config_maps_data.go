package configuration

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

type Data struct {
	autogen_bases.ChartSubcomponentParentClassTypes
	structs.SuperParentClass
}
