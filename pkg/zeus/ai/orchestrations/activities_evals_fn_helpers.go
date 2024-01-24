package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func TransformJSONToEvalScoredMetrics(jsonSchemaDef *artemis_orchestrations.JsonSchemaDefinition) (*artemis_orchestrations.EvalMetricsResults, error) {
	var metrics []artemis_orchestrations.EvalMetric
	for _, value := range jsonSchemaDef.Fields {
		for i, _ := range value.EvalMetrics {
			if value.EvalMetrics[i] == nil {
				value.EvalMetrics[i] = &artemis_orchestrations.EvalMetric{}
			}
			if value.EvalMetrics[i].EvalMetricResult == nil {
				chs := chronos.Chronos{}
				value.EvalMetrics[i].EvalMetricResult = &artemis_orchestrations.EvalMetricResult{
					EvalMetricResultID: aws.Int(chs.UnixTimeStampNow()),
				}
			}
			if value.EvalMetrics[i].EvalMetricComparisonValues == nil {
				value.EvalMetrics[i].EvalMetricComparisonValues = &artemis_orchestrations.EvalMetricComparisonValues{}
			}
			switch value.DataType {
			case "integer":
				if value.IntegerValue == nil {
					return nil, fmt.Errorf("no int value for key '%s'", value.FieldName)
				}
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil && value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
				}
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger != nil {
					value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetIntEvalComparisonResult(value.EvalMetrics[i].EvalOperator, *value.IntegerValue, aws.ToInt(value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger)))
				} else if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber != nil {
					value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetIntEvalComparisonResult(value.EvalMetrics[i].EvalOperator, *value.IntegerValue, int(aws.ToFloat64(value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))))
				} else {
					return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
				}
			case "number":
				if value.NumberValue == nil {
					return nil, fmt.Errorf("no number value for key '%s'", value.FieldName)
				}
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetNumericEvalComparisonResult(value.EvalMetrics[i].EvalOperator, *value.NumberValue, aws.ToFloat64(value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber)))
			case "string":
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString == nil {
					return nil, fmt.Errorf("no comparison string for key '%s'", value.FieldName)
				}
				if value.StringValue == nil {
					return nil, fmt.Errorf("no string value for key '%s'", value.FieldName)
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetStringEvalComparisonResult(value.EvalMetrics[i].EvalOperator, *value.StringValue, aws.ToString(value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString)))
			case "boolean":
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean == nil {
					return nil, fmt.Errorf("no comparison boolean for key '%s'", value.FieldName)
				}
				if value.BooleanValue == nil {
					return nil, fmt.Errorf("no boolean value for key '%s'", value.FieldName)
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetBooleanEvalComparisonResult(aws.ToBool(value.BooleanValue), aws.ToBool(value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean)))
			case "array[integer]":
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil && value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
				}
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger != nil {
					results, rerr := EvaluateIntArray(value.EvalMetrics[i].EvalOperator, value.IntegerValueSlice, *value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger)
					if rerr != nil {
						return nil, rerr
					}
					value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
				} else {
					results, rerr := EvaluateIntArray(value.EvalMetrics[i].EvalOperator, value.IntegerValueSlice, int(*value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
					if rerr != nil {
						return nil, rerr
					}
					value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
				}
			case "array[number]":
				if value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil {
					return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
				}
				results, rerr := EvaluateNumericArray(value.EvalMetrics[i].EvalOperator, value.NumberValueSlice, *value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber)
				if rerr != nil {
					return nil, rerr
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			case "array[string]":
				results, rerr := EvaluateStringArray(value.StringValueSlice, value.EvalMetrics[i].EvalOperator, *value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString)
				if rerr != nil {
					return nil, rerr
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			case "array[boolean]":
				results, rerr := EvaluateBooleanArray(value.BooleanValueSlice, *value.EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean)
				if rerr != nil {
					return nil, rerr
				}
				value.EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			default:
				return nil, fmt.Errorf("unknown data type '%s'", value.DataType)
			}
		}
	}
	return &artemis_orchestrations.EvalMetricsResults{
		EvalMetricsResults: metrics,
	}, nil
}

func GetBooleanEvalComparisonResult(actual, expected bool) bool {
	return actual == expected
}

func GetIntEvalComparisonResult(operator string, actual, expected int) bool {
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

func EvaluateIntArray(operator string, array []int, expected int) ([]bool, error) {
	var results []bool
	for _, value := range array {
		result := GetIntEvalComparisonResult(operator, value, expected)
		results = append(results, result)
	}
	return results, nil
}

func EvaluateNumericArray(operator string, array []float64, expected float64) ([]bool, error) {
	var results []bool
	for _, value := range array {
		result := GetNumericEvalComparisonResult(operator, value, expected)
		results = append(results, result)
	}
	return results, nil
}

func Pass(results []bool) bool {
	for _, result := range results {
		if !result {
			return false
		}
	}
	return true
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

func convertToInt(value interface{}) (int, error) {
	av, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("value is not float64")
	}
	return av, nil
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