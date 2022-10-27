package common_conversions

import "encoding/json"

func ParseParentChildMap(childElementMap map[string]interface{}, pc *PC) error {
	parentClassType, pClassTypeOK := childElementMap["chart_subcomponent_parent_class_types"]
	childClassType, cClassTypeOK := childElementMap["chart_subcomponent_child_class_types"]
	childValuesType, cClassValueOK := childElementMap["chart_subcomponents_child_values"]

	switch pClassTypeOK || cClassTypeOK || cClassValueOK {
	case pClassTypeOK:

		bytes, berr := getBytes(parentClassType)
		if berr != nil {
			return berr
		}
		perr := json.Unmarshal(bytes, &pc.ParentClassTypesDB)
		if perr != nil {
			return perr
		}
	case cClassTypeOK:
		bytes, berr := getBytes(childClassType)
		if berr != nil {
			return berr
		}
		perr := json.Unmarshal(bytes, &pc.ChartSubcomponentChildClassTypes)
		if perr != nil {
			return perr
		}
	case cClassValueOK:
		bytes, berr := getBytes(childValuesType)
		if berr != nil {
			return berr
		}
		perr := json.Unmarshal(bytes, &pc.ChartSubcomponentsChildValues)
		if perr != nil {
			return perr
		}
	}
	return nil
}

func getBytes(structToBytes interface{}) ([]byte, error) {
	bytes, berr := json.Marshal(structToBytes)
	return bytes, berr
}
