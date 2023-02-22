package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
	"github.com/zeus-fyi/olympus/pkg/utils/misc/dev_hacks"
)

func InsertStatefulSetDefinitions() error {
	cck := statefulsets.NewStatefulSet()

	err := dev_hacks.Use(cck)
	return err
}
