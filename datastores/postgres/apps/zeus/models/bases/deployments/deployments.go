package deployments

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
)

const ModelName = "Deployment"

type Deployment struct {
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
