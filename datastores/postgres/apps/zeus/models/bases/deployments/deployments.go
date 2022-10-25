package deployments

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
)

const ModelName = "Deployment"

type Deployment struct {
	KindDefinition autogen_bases.ChartComponentResources

	Metadata common.ParentMetaData
	Spec     Spec
}

type Spec struct {
	common.SpecWorkload
	Template containers.PodTemplateSpec
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.Metadata.Metadata = common.NewMetadata()
	d.Metadata.ChartSubcomponentParentClassTypeName = "deploymentParentMetadata"
	d.Spec = NewDeploymentSpec()
	d.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "deploymentSpec",
	}
	return d
}

func NewDeploymentSpec() Spec {
	ds := Spec{}
	ds.SpecWorkload = common.NewSpecWorkload()
	ds.Template = containers.NewPodTemplateSpec()
	return ds
}
