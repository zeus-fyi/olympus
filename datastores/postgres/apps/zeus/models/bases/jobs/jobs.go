package jobs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1Batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	K8sJob         *v1Batch.Job
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     v1Batch.JobSpec
}

const JobChartComponentResourceID = 0

func NewJob() Job {
	s := Job{}
	typeMeta := metav1.TypeMeta{
		Kind:       "Job",
		APIVersion: "batch/v1",
	}
	s.K8sJob = &v1Batch.Job{TypeMeta: typeMeta}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Job",
		ChartComponentApiVersion: "batch/v1",
		ChartComponentResourceID: JobChartComponentResourceID,
	}

	s.Metadata.Metadata = structs.NewMetadata()
	s.Metadata.ChartSubcomponentParentClassTypeName = "JobParentMetadata"
	s.Metadata.ChartComponentResourceID = JobChartComponentResourceID
	return s
}
