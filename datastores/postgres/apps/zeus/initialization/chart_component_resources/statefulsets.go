package chart_component_resources

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulset"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertStatefulSetDefinitions() error {
	cck := statefulset.NewStatefulSet()

	err := dev_hacks.Use(cck)
	return err
}
