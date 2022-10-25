package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertServiceDefinitions() error {
	cck := networking.NewService()

	err := dev_hacks.Use(cck)
	return err
}
