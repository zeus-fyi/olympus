package workloads

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/containers"
)

type Deployment struct {
	KindDefinition        autogen_structs.ChartComponentKinds
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	Spec     DeploymentSpec
}

type DeploymentSpec struct {
	Replicas int
	Selector common.Selector

	Template containers.PodTemplateSpec
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.KindDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.ParentClassDefinition = autogen_structs.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentKindID:                 0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "deploymentSpec",
	}
	d.Spec = NewDeploymentSpec()
	return d
}

func NewDeploymentSpec() DeploymentSpec {
	ds := DeploymentSpec{}
	ds.Selector = common.NewSelector()
	ds.Template = containers.NewPodTemplateSpec()
	return ds
}
