package chart_component_kinds

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func InsertStatefulSetDefinitions() error {
	cck := workloads.NewStatefulSet()

	err := dev_hacks.Use(cck)
	return err
}
