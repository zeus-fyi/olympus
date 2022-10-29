package configuration

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ConfigMapChartComponentResourceID = 12

type ConfigMap struct {
	K8sConfigMap   v1.ConfigMap
	KindDefinition autogen_bases.ChartComponentResources
	Metadata       structs.ParentMetaData

	Immutable *structs.ChildClassSingleValue

	// TODO give parent class names
	Data       Data
	BinaryData structs.SuperParentClass
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	typeMeta := metav1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	}
	cm.K8sConfigMap = v1.ConfigMap{
		TypeMeta:   typeMeta,
		ObjectMeta: metav1.ObjectMeta{},
		Immutable:  nil,
		Data:       nil,
		BinaryData: nil,
	}
	cm.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
		ChartComponentResourceID: ConfigMapChartComponentResourceID,
	}
	cm.Metadata.ChartComponentResourceID = ConfigMapChartComponentResourceID
	cm.Metadata.ChartSubcomponentParentClassTypeName = "ConfigMapParentMetadata"
	cm.Metadata.Metadata = structs.NewMetadata()

	return cm
}

func (cm *ConfigMap) SetChartPackageID(id int) {
	cm.Data.ChartPackageID = id
	cm.Metadata.ChartPackageID = id
}
