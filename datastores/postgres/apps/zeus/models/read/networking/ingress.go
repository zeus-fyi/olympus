package read_networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
)

func DBIngressResource(ing *ingresses.Ingress, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseDeploymentParentChildAggValues(ckagg)
	if pcerr != nil {
		return pcerr
	}

	err := ing.ParseDBConfigToK8s(pcGroupMap)
	return err
}
