package jobs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1Batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	K8sJob         *v1Batch.Job
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     JobSpec
}

type JobSpec struct {
	common.ParentClass
	structs.ChildClassSingleValue
}

const JobChartComponentResourceID = 6

func NewJob() Job {
	j := Job{}
	typeMeta := metav1.TypeMeta{
		Kind:       "Job",
		APIVersion: "batch/v1",
	}
	j.K8sJob = &v1Batch.Job{TypeMeta: typeMeta}
	j.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Job",
		ChartComponentApiVersion: "batch/v1",
		ChartComponentResourceID: JobChartComponentResourceID,
	}
	j.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             JobChartComponentResourceID,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "Spec",
	}
	j.Metadata.Metadata = structs.NewMetadata()
	j.Metadata.ChartSubcomponentParentClassTypeName = "JobParentMetadata"
	j.Metadata.ChartComponentResourceID = JobChartComponentResourceID
	return j
}
