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

func (r *Rules) AddIngressRule(hostName string, httpPathsSlice []string) {
	var ts chronos.Chronos
	rule := structs.NewSuperParentClassWithBothChildTypes("rules", "rules", "rules")
	childTypeID := ts.UnixTimeStampNow()

	// single values part
	rule.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, "rules", "host", hostName)

	// multi values part
	rule.AddKeyValuesAndUniqueChildID(childTypeID, "rules", "path", httpPathsSlice)
	r.SuperParentClassSlice = append(r.SuperParentClassSlice, rule)
}
