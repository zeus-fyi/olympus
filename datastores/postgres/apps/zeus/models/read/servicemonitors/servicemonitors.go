package read_servicemonitors

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/servicemonitors"
)

func DBServiceMonitorResource(sm *servicemonitors.ServiceMonitor, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		log.Err(pcerr)
		return pcerr
	}
	err := sm.ParseDBConfigToK8s(pcGroupMap)
	return err
}
