package workloads

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Deployment struct {
	KindDefinition autogen_bases.ChartComponentResources

	Metadata DeploymentMetadata
	Spec     Spec
}

type DeploymentMetadata struct {
	autogen_bases.ChartSubcomponentParentClassTypes
	common.Metadata
}

type Spec struct {
	autogen_bases.ChartSubcomponentParentClassTypes
	DeploymentSpec
}

type DeploymentSpec struct {
	Replicas common.ChildClassSingleValue
	Selector common.Selector

	Template containers.PodTemplateSpec
}

func (ds *DeploymentSpec) GetReplicaCount32IntPtr() *int32 {
	return string_utils.ConvertStringTo32BitPtrInt(ds.Replicas.ChartSubcomponentValue)
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	d.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
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
