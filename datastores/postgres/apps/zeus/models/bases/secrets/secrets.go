package secrets

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/core/v1"
)

const SecretChartComponentResourceID = 13

type Secret struct {
	K8sSecret      v1.Secret
	KindDefinition autogen_bases.ChartComponentResources
	Metadata       structs.ParentMetaData

	Immutable *structs.ChildClassSingleValue

	StringData StringData
	Type       SecretType
}

func NewSecret() Secret {
	s := Secret{
		K8sSecret:      v1.Secret{},
		KindDefinition: autogen_bases.ChartComponentResources{},
		Metadata:       structs.ParentMetaData{},
		Immutable:      nil,
		StringData:     StringData{},
		Type:           SecretType{},
	}

	return s
}
