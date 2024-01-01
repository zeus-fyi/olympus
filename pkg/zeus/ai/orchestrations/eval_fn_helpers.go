package ai_platform_service_orchestrations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
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

func TransformJSONToEvalScoredMetrics(dataMap map[string]interface{}, emMap map[string]artemis_orchestrations.EvalMetric) (*artemis_orchestrations.EvalMetricsResults, error) {
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
				av, aerr := convertToFloat64(value)
				if aerr != nil {
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
					return nil, aerr
				}
				result = GetNumericEvalComparisonResult(v.EvalOperator, av, *v.EvalComparisonNumber)
			case "string":
				if v.EvalComparisonString == nil {
					return nil, fmt.Errorf("no comparison string for key '%s'", key)
				}
				av, aerr := convertToString(value)
				if aerr != nil {
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
					return nil, aerr
				}
				result = GetStringEvalComparisonResult(v.EvalOperator, av, *v.EvalComparisonString)
			case "boolean":
				if v.EvalComparisonBoolean == nil {
					return nil, fmt.Errorf("no comparison boolean for key '%s'", key)
				}
				av, aerr := convertToBool(value)
				if aerr != nil {
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
					return nil, aerr
				}
				result = GetBooleanEvalComparisonResult(av, *v.EvalComparisonBoolean)
			case "array[number]":
				if v.EvalComparisonNumber == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", key)
				}
				ifs, rok := value.([]interface{})
				if !rok {
					aerr := fmt.Errorf("value is not []interface{} for key '%s'", key)
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: array[number] failed")
					return nil, aerr
				}
				sa, aerr := interfaceSliceToFloat64Slice(ifs)
				if aerr != nil {
					return nil, aerr
				}
				results, rerr := EvaluateNumericArray(v.EvalOperator, sa, *v.EvalComparisonNumber)
				if rerr != nil {
					return nil, rerr
				}
				result = ContainsFalse(results)
			case "array[string]":
				if v.EvalComparisonString == nil && !strings.Contains(v.EvalOperator, "all-unique-words") {
					return nil, fmt.Errorf("no comparison string for key '%s'", key)
				}
				ifs, rok := value.([]interface{})
				if !rok {
					aerr := fmt.Errorf("value is not []interface{} for key '%s'", key)
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: array[number] failed")
					return nil, aerr
				}
				sa, aerr := interfaceSliceToStringSlice(ifs)
				if aerr != nil {
					return nil, aerr
				}
				results, rerr := EvaluateStringArray(sa, v.EvalOperator, *v.EvalComparisonString)
				if rerr != nil {
					return nil, rerr
				}
				result = ContainsFalse(results)
			case "array[boolean]":
				ifs, rok := value.([]interface{})
				if !rok {
					aerr := fmt.Errorf("value is not []interface{} for key '%s'", key)
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: array[boolean] failed")
					return nil, aerr
				}
				ba, aerr := interfaceSliceToBoolSlice(ifs)
				if aerr != nil {
					log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: array[boolean] failed")
					return nil, aerr
				}
				results, rerr := EvaluateBooleanArray(ba, *v.EvalComparisonBoolean)
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
	return &artemis_orchestrations.EvalMetricsResults{
		EvalMetricsResults: metrics,
	}, nil
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

func EvaluateNumericArray(operator string, array []float64, expected float64) ([]bool, error) {
	var results []bool
	for _, value := range array {
		result := GetNumericEvalComparisonResult(operator, value, expected)
		results = append(results, result)
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

func EvaluateBooleanArray(array []bool, expected bool) ([]bool, error) {
	var results []bool
	for _, value := range array {
		result := GetBooleanEvalComparisonResult(value, expected)
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

func EvaluateStringArray(array []string, operator, expected string) ([]bool, error) {
	var results []bool
	seen := make(map[string]bool)
	for _, value := range array {
		if operator == "all-unique-words" {
			_, ok := seen[value]
			if ok {
				results = append(results, false)
			} else {
				results = append(results, true)
			}
		} else {
			result := GetStringEvalComparisonResult(operator, value, expected)
			results = append(results, result)
		}
		seen[value] = true
	}
	return results, nil
}

// Helper function to assert and convert value to string
func convertToString(value interface{}) (string, error) {
	av, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not string")
	}
	return av, nil
}
func interfaceSliceToFloat64Slice(interfaceSlice []interface{}) ([]float64, error) {
	float64Slice := make([]float64, len(interfaceSlice))
	for i, v := range interfaceSlice {
		f, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("value at index %d is not a float64", i)
		}
		float64Slice[i] = f
	}
	return float64Slice, nil
}
func interfaceSliceToBoolSlice(interfaceSlice []interface{}) ([]bool, error) {
	boolSlice := make([]bool, len(interfaceSlice))
	for i, v := range interfaceSlice {
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("value at index %d is not a bool", i)
		}
		boolSlice[i] = b
	}
	return boolSlice, nil
}
func interfaceSliceToStringSlice(interfaceSlice []interface{}) ([]string, error) {
	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value at index %d is not a string", i)
		}
		stringSlice[i] = str
	}
	return stringSlice, nil
}
func convertToFloat64(value interface{}) (float64, error) {
	av, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("value is not float64")
	}
	return av, nil
}

// Helper function to assert and convert value to bool
func convertToBool(value interface{}) (bool, error) {
	av, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("value is not bool")
	}
	return av, nil
}
