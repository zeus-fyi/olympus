package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func NewTLS() TLS {
	parentClassTypeName := "tls"

	tls := TLS{
		structs.NewSuperParentClassGroup(parentClassTypeName),
	}
	return tls
}

type TLS struct {
	structs.SuperParentClassGroup
}

func (t *TLS) AddIngressTLS(secretNameValue string, hosts []string) {
	var ts chronos.Chronos
	singleChildClassTypeName := "secretName"
	multiChildClassTypeName := "hosts"
	tlsIngress := structs.NewSuperParentClassWithBothChildTypes("ingress", singleChildClassTypeName, multiChildClassTypeName)

	childTypeID := ts.UnixTimeStampNow()
	tlsIngress.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, singleChildClassTypeName, singleChildClassTypeName, secretNameValue)
	tlsIngress.AddKeyValuesAndUniqueChildID(childTypeID, multiChildClassTypeName, "hosts", hosts)
	t.SuperParentClassSlice = append(t.SuperParentClassSlice, tlsIngress)

}
