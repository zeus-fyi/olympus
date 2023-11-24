package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (cj *CronJob) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := cj.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "CronJobParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&cj.K8sCronJob.ObjectMeta, pc)
		}
	}
	return nil
}

func (cj *CronJob) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "CronJobJobSpec":
			err := json.Unmarshal([]byte(value), &cj.K8sCronJob.Spec)
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}
	return nil
}
