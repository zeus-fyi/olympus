package chart_component_kinds

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertDeploymentDefinitions() error {
	cck := autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}

	err := dev_hacks.Use(cck)
	return err
}
