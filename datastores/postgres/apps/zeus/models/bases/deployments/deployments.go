package deployments

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ModelName = "Deployment"

type Deployment struct {
	K8sDeployment  *v1.Deployment
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     Spec
}

type Spec struct {
	structs.SpecWorkload
	Template containers.PodTemplateSpec
}

func NewDeployment() Deployment {
	d := Deployment{}
	typeMeta := metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	}
	d.K8sDeployment = &v1.Deployment{
		TypeMeta:   typeMeta,
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.DeploymentSpec{},
		Status:     v1.DeploymentStatus{},
	}
	d.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.Metadata.Metadata = structs.NewMetadata()
	d.Metadata.ChartSubcomponentParentClassTypeName = "DeploymentParentMetadata"
	d.Spec = NewDeploymentSpec()
	d.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "DeploymentSpec",
	}
	return d
}

func NewDeploymentSpec() Spec {
	ds := Spec{}
	ds.SpecWorkload = structs.NewSpecWorkload()
	ds.Template = containers.NewPodTemplateSpec()
	return ds
}
