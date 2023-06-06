package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	v1Batch "k8s.io/api/batch/v1"
)

func (j *Job) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := j.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "JobParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&j.K8sJob.ObjectMeta, pc)
		}
	}
	return nil
}

func (j *Job) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "JobSpec":
			j.K8sJob = v1Batch.Job{}
			err := json.Unmarshal([]byte(value), &j.K8sJob.Spec)
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}
	return nil
}
