package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func TransformJSONToEvalScoredMetrics(jsonSchemaDef *artemis_orchestrations.JsonSchemaDefinition) (*artemis_orchestrations.EvalMetricsResults, error) {
	var metrics []artemis_orchestrations.EvalMetricsResult
	var evmID int
	for _, value := range jsonSchemaDef.Fields {
		if value.EvalMetric == nil {
			continue
		}
		var result bool
		switch value.DataType {
		case "integer":
			if value.EvalMetric.EvalComparisonNumber == nil {
				return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
			}
			av, aerr := convertToInt(value)
			if aerr != nil {
				log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToInt failed")
				return nil, aerr
			}
			result = GetIntEvalComparisonResult(value.EvalMetric.EvalOperator, av, int(*value.EvalMetric.EvalComparisonNumber))
		case "number":
			if value.EvalMetric.EvalComparisonNumber == nil {
				return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
			}
			av, aerr := convertToFloat64(value)
			if aerr != nil {
				log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
				return nil, aerr
			}
			result = GetNumericEvalComparisonResult(value.EvalMetric.EvalOperator, av, *value.EvalMetric.EvalComparisonNumber)
		case "string":
			if value.EvalMetric.EvalComparisonString == nil {
				return nil, fmt.Errorf("no comparison string for key '%s'", value.FieldName)
			}
			av, aerr := convertToString(value)
			if aerr != nil {
				log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
				return nil, aerr
			}
			result = GetStringEvalComparisonResult(value.EvalMetric.EvalOperator, av, *value.EvalMetric.EvalComparisonString)
		case "boolean":
			if value.EvalMetric.EvalComparisonBoolean == nil {
				return nil, fmt.Errorf("no comparison boolean for key '%s'", value.FieldName)
			}
			av, aerr := convertToBool(value)
			if aerr != nil {
				log.Err(aerr).Msg("TransformJSONToEvalScoredMetrics: convertToFloat64 failed")
				return nil, aerr
			}
			result = GetBooleanEvalComparisonResult(av, *value.EvalMetric.EvalComparisonBoolean)
		case "array[integer]":
			if value.EvalMetric.EvalComparisonNumber == nil {
				return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
			}
			results, rerr := EvaluateIntArray(value.EvalMetric.EvalOperator, value.IntValueSlice, int(*value.EvalMetric.EvalComparisonNumber))
			if rerr != nil {
				return nil, rerr
			}
			result = Pass(results)
		case "array[number]":
			if value.EvalMetric.EvalComparisonNumber == nil {
				return nil, fmt.Errorf("no comparison number for key '%s'", value.FieldName)
			}
			results, rerr := EvaluateNumericArray(value.EvalMetric.EvalOperator, value.NumberValueSlice, *value.EvalMetric.EvalComparisonNumber)
			if rerr != nil {
				return nil, rerr
			}
			result = Pass(results)
		case "array[string]":
			results, rerr := EvaluateStringArray(value.StringValueSlice, value.EvalMetric.EvalOperator, *value.EvalMetric.EvalComparisonString)
			if rerr != nil {
				return nil, rerr
			}
			result = Pass(results)
		case "array[boolean]":
			results, rerr := EvaluateBooleanArray(value.BooleanValueSlice, *value.EvalMetric.EvalComparisonBoolean)
			if rerr != nil {
				return nil, rerr
			}
			result = Pass(results)
		}
		if value.EvalMetric.EvalMetricID != nil {
			evmID = *value.EvalMetric.EvalMetricID
		}
		metrics = append(metrics, artemis_orchestrations.EvalMetricsResult{
			EvalMetricID:      evmID,
			EvalResultOutcome: result,
		})

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
