package ai_platform_service_orchestrations

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func GetBooleanEvalComparisonResult(actual, expected bool) bool {
	return actual == expected
}

func GetIntEvalComparisonResult(operator string, actual, expected int) bool {
	switch operator {
	case "==", "eq":
		return actual == expected
	case "!=", "neq":
		return actual != expected
	case ">", "gt":
		return actual > expected
	case "<", "lt":
		return actual < expected
	case ">=", "gte":
		return actual >= expected
	case "<=", "lte":
		return actual <= expected
	}
	return false
}

func GetNumericEvalComparisonResult(operator string, actual, expected float64) bool {
	switch operator {
	case "==", "eq":
		return actual == expected
	case "!=", "neq":
		return actual != expected
	case ">", "gt":
		return actual > expected
	case "<", "lt":
		return actual < expected
	case ">=", "gte":
		return actual >= expected
	case "<=", "lte":
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
	case "equals-one-from-list":
		acceptable := strings.Split(expected, ",")
		for _, a := range acceptable {
			if strings.TrimSpace(actual) == strings.TrimSpace(a) {
				return true
			}
		}
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
	case "length-eq":
		expectedLen := len(expected)
		comparedLengthLimit, err := strconv.Atoi(expected)
		if err == nil {
			expectedLen = comparedLengthLimit
		}
		actualLen := len(actual)
		return actualLen == expectedLen
	case "length-less-than":
		expectedLen := len(expected)
		comparedLengthLimit, err := strconv.Atoi(expected)
		if err == nil {
			expectedLen = comparedLengthLimit
		}
		actualLen := len(actual)
		if actualLen < expectedLen {
			return true
		}
	case "length-less-than-eq":
		expectedLen := len(expected)
		comparedLengthLimit, err := strconv.Atoi(expected)
		if err == nil {
			expectedLen = comparedLengthLimit
		}
		if len(actual) <= expectedLen {
			return true
		}
	case "length-greater-than":
		expectedLen := len(expected)
		comparedLengthLimit, err := strconv.Atoi(expected)
		if err == nil {
			expectedLen = comparedLengthLimit
		}
		if len(actual) > expectedLen {
			return true
		}
	case "length-greater-than-eq":
		expectedLen := len(expected)
		comparedLengthLimit, err := strconv.Atoi(expected)
		if err == nil {
			expectedLen = comparedLengthLimit
		}
		if len(actual) >= expectedLen {
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

func haveIdenticalKeys[K comparable, V any](map1, map2 map[K]*V) bool {
	if map1 == nil || map2 == nil {
		return false // Handle nil maps
	}
	for key := range map1 {
		if _, exists := map2[key]; !exists {
			return false
		}
	}
	return true
}

func copyMatchingFieldValues(tasksSchemaMap, schemasMap map[string]*artemis_orchestrations.JsonSchemaDefinition) {
	if tasksSchemaMap == nil || schemasMap == nil {
		return // Handle nil maps
	}
	if !haveIdenticalKeys(tasksSchemaMap, schemasMap) {
		return // The maps do not have identical keys
	}
	// Ensure FieldsMap is populated for both maps
	for si, _ := range tasksSchemaMap {
		populateFieldsMap(tasksSchemaMap[si])
	}
	for si, _ := range schemasMap {
		populateFieldsMap(schemasMap[si])
	}
	for schemaID, _ := range tasksSchemaMap {
		copyFieldValues(tasksSchemaMap[schemaID], schemasMap[schemaID])
	}
}

func copyMatchingFieldValuesFromResp(respJsonResults *artemis_orchestrations.JsonSchemaDefinition, schemasMap map[string]*artemis_orchestrations.JsonSchemaDefinition) bool {
	if schemasMap == nil || respJsonResults == nil {
		return false // Handle nil maps
	}
	evalSchema, ok := schemasMap[respJsonResults.SchemaStrID]
	if !ok || evalSchema == nil {
		return false // The maps do not have identical keys
	}
	populateFieldsMap(respJsonResults)
	for _, schema := range schemasMap {
		populateFieldsMap(schema)
	}
	return copyFieldValues(respJsonResults, evalSchema)
}

func copyFieldValues(src, dest *artemis_orchestrations.JsonSchemaDefinition) bool {
	if src == nil || dest == nil {
		return false // Handle nil arguments
	}

	for _, srcField := range src.FieldsMap {
		if srcField == nil {
			continue // Skip if srcField is nil
		}

		destField, ok := dest.FieldsMap[srcField.FieldStrID]
		if !ok || destField == nil {
			return false
		}
		if !srcField.FieldValue.IsValidated {
			return false
		}

		if srcField.DataType == destField.DataType {
			destField.FieldValue = srcField.FieldValue
		}
	}
	return true
}

func populateFieldsMap(schema *artemis_orchestrations.JsonSchemaDefinition) {
	if schema == nil || schema.Fields == nil {
		return
	}
	schema.FieldsMap = make(map[string]*artemis_orchestrations.JsonSchemaField)
	for i, _ := range schema.Fields {
		field := &schema.Fields[i]
		schema.FieldsMap[field.FieldStrID] = field
	}
}
