package read_configuration

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
)

func DBConfigMapResource(cm *configuration.ConfigMap, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		return pcerr
	}

	err := cm.ParseDBConfigToK8s(pcGroupMap)
	return err
}
