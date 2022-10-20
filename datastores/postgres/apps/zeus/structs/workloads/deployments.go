package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

type Deployment struct {
	KindDefinition autogen_structs.ChartComponentResources

	Metadata DeploymentMetadata
	Spec     Spec
}

type DeploymentMetadata struct {
	autogen_structs.ChartSubcomponentParentClassTypes
	common.Metadata
}

type Spec struct {
	autogen_structs.ChartSubcomponentParentClassTypes
	DeploymentSpec
}

type DeploymentSpec struct {
	Replicas common.ChildClassSingleValue
	Selector common.Selector

	Template containers.PodTemplateSpec
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.KindDefinition = autogen_structs.ChartComponentResources{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.Spec.ChartSubcomponentParentClassTypes = autogen_structs.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "deploymentSpec",
	}
	d.Metadata.Metadata = common.NewMetadata()
	d.Spec.DeploymentSpec = NewDeploymentSpec()
	return d
}

func NewDeploymentSpec() DeploymentSpec {
	ds := DeploymentSpec{}
	ds.Selector = common.NewSelector()
	ds.Template = containers.NewPodTemplateSpec()
	ds.Replicas = common.NewInitChildClassSingleValue("replicas", "0")
	return ds
}
