package jobs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1Batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const CronJobChartComponentResourceID = 7

type CronJob struct {
	K8sCronJob     v1Batch.CronJob
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     CronJobSpec
}

type CronJobSpec struct {
	common.ParentClass
	structs.ChildClassSingleValue
}

func NewCronJob() CronJob {
	cj := CronJob{}
	typeMeta := metav1.TypeMeta{
		Kind:       "CronJob",
		APIVersion: "batch/v1",
	}
	cj.K8sCronJob = v1Batch.CronJob{TypeMeta: typeMeta}
	cj.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "CronJob",
		ChartComponentApiVersion: "batch/v1",
		ChartComponentResourceID: CronJobChartComponentResourceID,
	}
	cj.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             CronJobChartComponentResourceID,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "Spec",
	}
	cj.Metadata.Metadata = structs.NewMetadata()
	cj.Metadata.ChartSubcomponentParentClassTypeName = "CronJobParentMetadata"
	cj.Metadata.ChartComponentResourceID = CronJobChartComponentResourceID
	return cj
}
