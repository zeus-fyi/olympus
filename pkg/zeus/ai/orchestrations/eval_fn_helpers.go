package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
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

func TransformJSONToEvalScoredMetrics(jsonData string, emMap map[string]artemis_orchestrations.EvalMetric) (*artemis_orchestrations.EvalMetricsResults, error) {
	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &dataMap); err != nil {
		return nil, err
	}
	var metrics []artemis_orchestrations.EvalMetricsResult
	var evmID int
	for key, value := range dataMap {
		dataType, err := getDataType(value)
		if err != nil {
			return nil, fmt.Errorf("error determining data type for key '%s': %v", key, err)
		}

		var result bool
		if v, ok := emMap[key]; ok {
			switch dataType {
			case "number":
				if v.EvalComparisonNumber == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", key)
				}
				av, fok := value.(float64)
				if !fok {
					return nil, fmt.Errorf("value is not float64 for key '%s'", key)
				}
				result = GetNumericEvalComparisonResult(v.EvalOperator, av, *v.EvalComparisonNumber)
			case "string":
				if v.EvalComparisonString == nil {
					return nil, fmt.Errorf("no comparison string for key '%s'", key)
				}
				av, fok := value.(string)
				if !fok {
					return nil, fmt.Errorf("value is not string for key '%s'", key)
				}
				result = GetStringEvalComparisonResult(v.EvalOperator, av, *v.EvalComparisonString)
			case "boolean":
				if v.EvalComparisonBoolean == nil {
					return nil, fmt.Errorf("no comparison boolean for key '%s'", key)
				}
				av, fok := value.(bool)
				if !fok {
					return nil, fmt.Errorf("value is not bool for key '%s'", key)
				}
				result = GetBooleanEvalComparisonResult(av, *v.EvalComparisonBoolean)
			case "array[number]":
				if v.EvalComparisonNumber == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", key)
				}
				results, rerr := EvaluateNumericArray(v.EvalOperator, value, *v.EvalComparisonNumber)
				if rerr != nil {
					return nil, rerr
				}
				result = ContainsFalse(results)
			case "array[string]":
				if v.EvalComparisonString == nil {
					return nil, fmt.Errorf("no comparison string for key '%s'", key)
				}
				results, rerr := EvaluateStringArray(value, v.EvalOperator, *v.EvalComparisonString)
				if rerr != nil {
					return nil, rerr
				}
				result = ContainsFalse(results)
			case "array[boolean]":
				results, rerr := EvaluateBooleanArray(value, *v.EvalComparisonBoolean)
				if rerr != nil {
					return nil, rerr
				}
				result = ContainsFalse(results)
			}
			if v.EvalMetricID != nil {
				evmID = *v.EvalMetricID
			}
			metrics = append(metrics, artemis_orchestrations.EvalMetricsResult{
				EvalMetricID:      evmID,
				EvalResultOutcome: result,
			})
		}
	}
	return nil, nil
}

func GetBooleanEvalComparisonResult(actual, expected bool) bool {
	return actual == expected
}

func GetNumericEvalComparisonResult(operator string, actual, expected float64) bool {
	switch operator {
	case "==":
		return actual == expected
	case "!=":
		return actual != expected
	case ">":
		return actual > expected
	case "<":
		return actual < expected
	case ">=":
		return actual >= expected
	case "<=":
		return actual <= expected
	}
	return false
}

func EvaluateNumericArray(operator string, array interface{}, expected float64) ([]bool, error) {
	var results []bool

	// Use reflection to handle different numeric array types
	rv := reflect.ValueOf(array)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice, got %s", rv.Kind())
	}

	for i := 0; i < rv.Len(); i++ {
		value := rv.Index(i)

		// Handle different types of numeric values
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result := GetNumericEvalComparisonResult(operator, float64(value.Int()), expected)
			results = append(results, result)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result := GetNumericEvalComparisonResult(operator, float64(value.Uint()), expected)
			results = append(results, result)
		case reflect.Float32, reflect.Float64:
			result := GetNumericEvalComparisonResult(operator, value.Float(), expected)
			results = append(results, result)
		default:
			return nil, fmt.Errorf("unsupported slice element type %s", value.Kind())
		}
	}

	return results, nil
}

func ContainsFalse(results []bool) bool {
	for _, result := range results {
		if !result {
			return true
		}
	}
	return false
}

func EvaluateBooleanArray(array interface{}, expected bool) ([]bool, error) {
	var results []bool

	rv := reflect.ValueOf(array)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice, got %s", rv.Kind())
	}

	for i := 0; i < rv.Len(); i++ {
		value := rv.Index(i)

		if value.Kind() != reflect.Bool {
			return nil, fmt.Errorf("expected a boolean slice, got slice of %s", value.Kind())
		}

		result := GetBooleanEvalComparisonResult(value.Bool(), expected)
		results = append(results, result)
	}

	return results, nil
}

func GetStringEvalComparisonResult(operator string, actual, expected string) bool {
	switch operator {
	case "contains":
		if strings.Contains(actual, expected) {
			return true
		}
	case "has-prefix":
		if strings.HasPrefix(actual, expected) {
			return true
		}
	case "has-suffix":
		if strings.HasSuffix(actual, expected) {
			return true
		}
	case "does-not-start-with-any":
		fs := &strings_filter.FilterOpts{
			DoesNotStartWithThese: strings.Split(expected, ","),
		}
		return strings_filter.FilterStringWithOpts(actual, fs)
	case "does-not-include":
		fs := &strings_filter.FilterOpts{
			DoesNotInclude: strings.Split(expected, ","),
		}
		return strings_filter.FilterStringWithOpts(actual, fs)
	case "equals":
		return actual == expected
	case "length-less-than":
		if len(actual) < len(expected) {
			return true
		}
	case "length-less-than-eq":
		if len(actual) <= len(expected) {
			return true
		}
	case "length-greater-than":
		if len(actual) > len(expected) {
			return true
		}
	case "length-greater-than-eq":
		if len(actual) >= len(expected) {
			return true
		}
	}
	return false
}
func EvaluateStringArray(array interface{}, operator, expected string) ([]bool, error) {
	var results []bool

	rv := reflect.ValueOf(array)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice, got %s", rv.Kind())
	}

	for i := 0; i < rv.Len(); i++ {
		value := rv.Index(i)

		if value.Kind() != reflect.String {
			return nil, fmt.Errorf("expected a string slice, got slice of %s", value.Kind())
		}

		actual := value.String()
		result := GetStringEvalComparisonResult(operator, actual, expected)
		results = append(results, result)
	}

	return results, nil
}
