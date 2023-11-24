package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
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
			ConvertMetadata(&j.K8sJob, pc)
		}
	}
	return nil
}

func ConvertMetadata(j *v1Batch.Job, pcSlice []common_conversions.PC) {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "labels":
			if j.Labels == nil {
				j.Labels = make(map[string]string)
			}
			j.Labels[pc.ChartSubcomponentKeyName] = pc.ChartSubcomponentValue
		case "annotations":
			if j.Annotations == nil {
				j.Annotations = make(map[string]string)
			}
			j.Annotations[pc.ChartSubcomponentKeyName] = pc.ChartSubcomponentValue
		case "name":
			j.Name = pc.ChartSubcomponentValue
		}
	}
}
func (j *Job) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "JobSpec":
			err := json.Unmarshal([]byte(value), &j.K8sJob.Spec)
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}
	return nil
}
