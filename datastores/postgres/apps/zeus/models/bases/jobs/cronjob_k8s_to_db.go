package jobs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func (cj *CronJob) ConvertK8CronJobToDB() error {
	cj.Metadata.ChartSubcomponentParentClassTypeName = "CronJobParentMetadata"
	cj.Metadata.Metadata = common_conversions.ConvertMetadata(cj.K8sCronJob.ObjectMeta)
	cj.Metadata.ChartComponentResourceID = CronJobChartComponentResourceID
	err := cj.ConvertK8sCronJobSpecToDB()
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

func (cj *CronJob) ConvertK8sCronJobSpecToDB() error {
	bytes, err := json.Marshal(cj.K8sCronJob.Spec)
	if err != nil {
		log.Err(err)
		return err
	}
	cj.Spec.ChildClassSingleValue = structs.NewInitChildClassSingleValue("spec", string(bytes))
	cj.Spec.ChartSubcomponentParentClassTypeName = "Spec"
	cj.Spec.ChildClassSingleValue.ChartSubcomponentChildClassTypeName = "CronJobSpec"
	return nil
}
