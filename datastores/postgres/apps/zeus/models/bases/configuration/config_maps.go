package configuration

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/core/v1"
)

const ConfigMapChartComponentResourceID = 12

type ConfigMap struct {
	K8sConfigMap   v1.ConfigMap
	KindDefinition autogen_bases.ChartComponentResources
	Metadata       structs.ParentMetaData

	Immutable *structs.ChildClassSingleValue

	// TODO give parent class names
	Data       structs.SuperParentClass
	BinaryData structs.SuperParentClass
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.K8sConfigMap = v1.ConfigMap{}
	cm.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
		ChartComponentResourceID: ConfigMapChartComponentResourceID,
	}
	cm.Metadata.ChartComponentResourceID = ConfigMapChartComponentResourceID
	cm.Metadata.ChartSubcomponentParentClassTypeName = "ConfigMapParentMetadata"
	cm.Metadata.Metadata = structs.NewMetadata()

	// TODO add data type
	return cm
}
