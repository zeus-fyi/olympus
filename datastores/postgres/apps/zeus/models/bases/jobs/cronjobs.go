package jobs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1Batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const CronJobChartComponentResourceID = 0

type CronJob struct {
	K8sCronJob     *v1Batch.CronJob
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     v1Batch.CronJobSpec
}

func NewCronJob() CronJob {
	s := CronJob{}
	typeMeta := metav1.TypeMeta{
		Kind:       "CronJob",
		APIVersion: "batch/v1",
	}
	s.K8sCronJob = &v1Batch.CronJob{TypeMeta: typeMeta}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "CronJob",
		ChartComponentApiVersion: "batch/v1",
		ChartComponentResourceID: CronJobChartComponentResourceID,
	}

	s.Metadata.Metadata = structs.NewMetadata()
	s.Metadata.ChartSubcomponentParentClassTypeName = "CronJobParentMetadata"
	s.Metadata.ChartComponentResourceID = CronJobChartComponentResourceID
	return s
}
