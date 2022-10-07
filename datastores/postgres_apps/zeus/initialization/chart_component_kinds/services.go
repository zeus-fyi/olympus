package chart_component_kinds

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertServiceDefinitions() error {
	cck := networking.NewService()

	err := dev_hacks.Use(cck)
	return err
}
