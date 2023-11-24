package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func (j *Job) ConvertK8JobToDB() error {
	jo := NewJob()
	j.KindDefinition = jo.KindDefinition
	j.Metadata = jo.Metadata
	j.Spec = jo.Spec
	j.Metadata.Name.ChartSubcomponentValue = j.K8sJob.Name
	j.Metadata.Metadata = common_conversions.CreateMetadataByFields(j.K8sJob.Name, j.K8sJob.Annotations, j.K8sJob.Labels)
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
