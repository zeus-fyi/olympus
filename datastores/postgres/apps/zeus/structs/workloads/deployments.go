package workloads

import (
	autogen_structs2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	common2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

type Deployment struct {
	KindDefinition        autogen_structs2.ChartComponentKinds
	ParentClassDefinition autogen_structs2.ChartSubcomponentParentClassTypes

	Metadata common2.Metadata
	Spec     DeploymentSpec
}

type DeploymentSpec struct {
	Replicas int
	Selector common2.Selector

	Template containers.PodTemplateSpec
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.KindDefinition = autogen_structs2.ChartComponentKinds{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.ParentClassDefinition = autogen_structs2.ChartSubcomponentParentClassTypes{
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
	ds.Selector = common2.NewSelector()
	ds.Template = containers.NewPodTemplateSpec()
	return ds
}
