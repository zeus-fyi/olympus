package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func (j *Job) ConvertK8JobToDB() error {
	j.Metadata.ChartSubcomponentParentClassTypeName = "JobParentMetadata"
	j.Metadata.Metadata = common_conversions.ConvertMetadata(j.K8sJob.ObjectMeta)
	j.Metadata.ChartComponentResourceID = JobChartComponentResourceID
	err := j.ConvertK8sJobSpecToDB()
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

func (j *Job) ConvertK8sJobSpecToDB() error {
	bytes, err := json.Marshal(j.K8sJob.Spec)
	if err != nil {
		log.Err(err)
		return err
	}
	j.Spec.ChildClassSingleValue = structs.NewInitChildClassSingleValue("spec", string(bytes))
	j.Spec.ChartSubcomponentParentClassTypeName = "Spec"
	j.Spec.ChildClassSingleValue.ChartSubcomponentChildClassTypeName = "JobSpec"
	return nil
}
