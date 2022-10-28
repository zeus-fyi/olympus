package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

func NewRule() Rule {
	parentClassTypeName := "rules"
	singleChildClassTypeName := "secretName"
	multiChildClassTypeName := "hosts"
	rule := Rule{
		structs.NewSuperParentClassWithBothChildTypes(parentClassTypeName, singleChildClassTypeName, multiChildClassTypeName),
	}
	return rule
}

type Rule struct {
	structs.SuperParentClass
}

//func (t *TLS) AddSecretName(secretNameValue string) {
//	t.ChildClassSingleValue.SetKeyAndValue("secretName", secretNameValue)
//}
//
//func (t *TLS) AddHosts(hostNamesMap map[string]string) {
//	t.ChildClassMultiValue.AddValues(hostNamesMap)
//}
