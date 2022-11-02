package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/pkg/utils/misc/dev_hacks"
)

func InsertServiceDefinitions() error {
	cck := services.NewService()

	err := dev_hacks.Use(cck)
	return err
}
