package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func NewRules() Rules {
	parentClassTypeName := "rules"
	rule := Rules{
		structs.NewSuperParentClassGroup(parentClassTypeName),
	}
	return rule
}

type Rules struct {
	structs.SuperParentClassGroup
}

func (r *Rules) AddNewIngressRuleAndUniqueChildClassID(hostName string, httpPathsMap map[string]string) {
	var ts chronos.Chronos
	rule := structs.NewSuperParentClassWithBothChildTypes("rules", "host", "http")
	childTypeID := ts.UnixTimeStampNow()
	rule.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, "host", "host", hostName)
	rule.AddValuesAndUniqueChildID(childTypeID, httpPathsMap)
	r.SuperParentClassSlice = append(r.SuperParentClassSlice, rule)
}
