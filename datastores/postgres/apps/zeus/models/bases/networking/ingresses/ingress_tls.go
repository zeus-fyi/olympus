package ingresses

import (
	"fmt"

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
	childTypeID := ts.UnixTimeStampNow()
	tlsGroupID := fmt.Sprintf("tls_%d", childTypeID)

	tlsIngress := structs.NewSuperParentClassWithBothChildTypes("tls", tlsGroupID, tlsGroupID)
	tlsIngress.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, tlsGroupID, "secretName", secretNameValue)
	tlsIngress.AddKeyValuesAndUniqueChildID(childTypeID, tlsGroupID, "hosts", hosts)
	t.SuperParentClassSlice = append(t.SuperParentClassSlice, tlsIngress)
}
