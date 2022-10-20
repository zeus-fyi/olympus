package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/configuration"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertConfigMapDefinitions() error {
	cck := configuration.NewConfigMap()

	err := dev_hacks.Use(cck)
	return err
}
