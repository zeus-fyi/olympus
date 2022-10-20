package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertDeploymentDefinitions() error {
	cck := workloads.NewDeployment()

	err := dev_hacks.Use(cck)
	return err
}
