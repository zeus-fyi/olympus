package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/pkg/utils/misc/dev_hacks"
)

func InsertConfigMapDefinitions() error {
	cck := configuration.NewConfigMap()

	err := dev_hacks.Use(cck)
	return err
}
