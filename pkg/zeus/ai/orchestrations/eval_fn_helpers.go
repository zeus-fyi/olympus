package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

// mapDataType maps the custom data type strings to jsonschema.DataType and handles array types.
func mapDataType(customType string) (jsonschema.DataType, *jsonschema.Definition, error) {
	switch customType {
	case "number":
		return jsonschema.Number, nil, nil
	case "string":
		return jsonschema.String, nil, nil
	case "boolean":
		return jsonschema.Boolean, nil, nil
	case "array[number]":
		return jsonschema.Array, &jsonschema.Definition{Type: jsonschema.Number}, nil
	case "array[string]":
		return jsonschema.Array, &jsonschema.Definition{Type: jsonschema.String}, nil
	case "array[boolean]":
		return jsonschema.Array, &jsonschema.Definition{Type: jsonschema.Boolean}, nil
	default:
		return "", nil, fmt.Errorf("unsupported data type: %s", customType)
	}
}

func getDataType(value interface{}) (string, error) {
	valueType := reflect.TypeOf(value)
	if valueType == nil {
		return "", fmt.Errorf("null values are not supported")
	}

	switch valueType.Kind() {
	case reflect.Float64:
		return "number", nil // JSON numbers are always floats
	case reflect.String:
		return "string", nil
	case reflect.Bool:
		return "boolean", nil
	case reflect.Slice, reflect.Array:
		// We need to check the type of the elements in the array.
		sliceValue := reflect.ValueOf(value)
		if sliceValue.Len() == 0 {
			return "array", nil // Empty array, cannot determine type of elements
		}
		firstElemType, err := getDataType(sliceValue.Index(0).Interface())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("array[%s]", firstElemType), nil
	default:
		return "", fmt.Errorf("unsupported data type: %v", valueType.Kind())
	}
}

func TransformJSONToEvalScoredMetrics(jsonData string, emMap map[string]artemis_orchestrations.EvalMetric) ([]artemis_orchestrations.EvalMetric, error) {
	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &dataMap); err != nil {
		return nil, err
	}

	var metrics []artemis_orchestrations.EvalMetric
	for key, value := range dataMap {
		dataType, err := getDataType(value)
		if err != nil {
			return nil, fmt.Errorf("error determining data type for key '%s': %v", key, err)
		}

		if _, ok := emMap[key]; ok {
			fmt.Println(value)
		}

		metrics = append(metrics, artemis_orchestrations.EvalMetric{
			EvalMetricName:     key,
			EvalMetricDataType: dataType,
		})
	}

	return metrics, nil
}
