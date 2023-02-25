package common_conversions

import (
	"encoding/json"
)

func ParseParentChildAggValues(ckaggString string) (ParentChildDB, error) {
	pcGroupMap := NewPCGroupMap()
	// represents the json compressed version in db that wraps parent, child, and values
	m := make(map[string][][]map[string][]map[string]interface{})
	err := json.Unmarshal([]byte(ckaggString), &m)
	if err != nil {
		return pcGroupMap, err
	}

	for resourceKind, parentChildContainersSlice := range m {
		if resourceKind == "parentWrapper" {
		}
		// resourceKind e.g. Deployment
		// parentChildContainersSlice is a list of all the parent element types of one parent eg DeploymentParentMetadata & Spec
		for _, singleParentTypeKey := range parentChildContainersSlice {
			// parentChildContainersSlice is a list of all the parent element types of one parent eg DeploymentParentMetadata
			for _, parentTypeNameMap := range singleParentTypeKey {
				// key is element name eg DeploymentParentMetadata
				for parentTypeName, parentChildObjSlice := range parentTypeNameMap {
					// _ = parentTypeName; = key for map
					pc := PC{}

					for _, childElementMap := range parentChildObjSlice {
						// child element types
						perr := ParseParentChildMap(childElementMap, &pc)
						if perr != nil {
							return pcGroupMap, perr
						}
					}
					pcGroupMap.AppendPCElementToMap(parentTypeName, pc)
				}
			}
		}
	}
	return pcGroupMap, nil
}
