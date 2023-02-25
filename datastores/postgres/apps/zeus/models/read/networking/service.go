package read_networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
)

func DBServiceResource(svc *services.Service, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		return pcerr
	}

	err := svc.ParsePCGroupMap(pcGroupMap)
	return err
}
