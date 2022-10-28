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
	rule := structs.NewSuperParentClassWithBothChildTypes("rules", "host", "http")
	childTypeID := ts.UnixTimeStampNow()

	// single values part
	rule.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, "host", "host", hostName)

	// multi values part
	rule.AddKeyValuesAndUniqueChildID(childTypeID, "http", "path", httpPathsSlice)
	r.SuperParentClassSlice = append(r.SuperParentClassSlice, rule)
}
