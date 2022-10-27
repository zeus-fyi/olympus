package read_deployments

import (
	"encoding/json"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Tmp struct {
	PCWrapper map[string][][]PC `db:"Deployment"`
}

type ParentClassTypesTMP struct {
	ChartSubcomponentParentClassTypeID   int    `db:"chart_subcomponent_parent_class_type_id"`
	ChartSubcomponentParentClassTypeName string `db:"chart_subcomponent_parent_class_type_name"`
}
type PC struct {
	ParentClassTypesTMP                            `db:"chart_subcomponent_parent_class_types"` // chart_subcomponent_parent_class_types
	autogen_bases.ChartSubcomponentChildClassTypes `db:"chart_subcomponent_child_class_types"`  //chart_subcomponent_child_class_types
	autogen_bases.ChartSubcomponentsChildValues    `db:"chart_subcomponents_child_values"`      // chart_subcomponents_child_values
}

func ParseDeploymentParentChildAggValues(ckaggString string) error {
	m := make(map[string][][]map[string][]map[string]interface{})
	err := json.Unmarshal([]byte(ckaggString), &m)
	if err != nil {
		return err
	}
	for resourceKind, parentChildContainersSlice := range m {
		if resourceKind == "parentWrapper" {
			//switch xx {
			//case "DeploymentParentMetadata":
			//	//_, _ = ParseMetadataValues(bytesN)
			//}
		}
		// resourceKind eg Deployment
		// parentChildContainersSlice is a list of all the parent element types of one parent eg DeploymentParentMetadata & Spec
		for _, singleParentTypeKey := range parentChildContainersSlice {
			// parentChildContainersSlice is a list of all the parent element types of one parent eg DeploymentParentMetadata
			for _, parentTypeNameMap := range singleParentTypeKey {
				// key is element name eg DeploymentParentMetadata
				for parentTypeName, parentChildObjSlice := range parentTypeNameMap {
					// _ = parentTypeName; = key for map
					dev_hacks.Use(parentTypeName)
					for _, childElementMap := range parentChildObjSlice {
						// child element types
						if _, ok := childElementMap["chart_subcomponent_parent_class_types"]; ok {
							dev_hacks.Use(childElementMap)
						}
						if _, ok := childElementMap["chart_subcomponent_child_class_types"]; ok {
							dev_hacks.Use(childElementMap)
						}
						if _, ok := childElementMap["chart_subcomponents_child_values"]; ok {
							dev_hacks.Use(childElementMap)
						}

					}

				}
			}

		}

	}
	return nil
}

func ParseMetadataValues(metadataStringBytes []byte) (metav1.ObjectMeta, error) {
	metaData := metav1.ObjectMeta{}
	err := json.Unmarshal(metadataStringBytes, &metaData)
	if err != nil {
		return metaData, err
	}
	return metaData, nil
}
