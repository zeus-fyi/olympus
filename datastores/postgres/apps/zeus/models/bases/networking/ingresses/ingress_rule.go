package ingresses

import (
	"fmt"

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
	childTypeID := ts.UnixTimeStampNow()
	rulesGroupID := fmt.Sprintf("rules_%d", childTypeID)
	rule := structs.NewSuperParentClassWithBothChildTypes("rules", rulesGroupID, rulesGroupID)

	// single values part
	rule.SetSingleChildClassIDTypeNameKeyAndValue(childTypeID, rulesGroupID, "host", hostName)

	// multi values part
	rule.AddKeyValuesAndUniqueChildID(childTypeID, rulesGroupID, "http", httpPathsSlice)
	r.SuperParentClassSlice = append(r.SuperParentClassSlice, rule)
}
