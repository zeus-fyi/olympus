package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func NewTLS() TLS {
	parentClassTypeName := "tls"
	singleChildClassTypeName := "secretName"
	multiChildClassTypeName := "hosts"
	tls := TLS{
		structs.NewSuperParentClassWithBothChildTypes(parentClassTypeName, singleChildClassTypeName, multiChildClassTypeName),
	}
	return tls
}

type TLS struct {
	structs.SuperParentClass
}

func (t *TLS) AddSecretName(secretNameValue string) {
	t.ChildClassSingleValue.SetKeyAndValue("secretName", secretNameValue)
}

func (t *TLS) AddHosts(hostNamesMap map[string]string) {
	t.ChildClassMultiValue.AddValues(hostNamesMap)
}
