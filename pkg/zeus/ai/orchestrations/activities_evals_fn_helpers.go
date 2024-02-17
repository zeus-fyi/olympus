package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

func TransformJSONToEvalScoredMetrics(jsonSchemaDef *artemis_orchestrations.JsonSchemaDefinition) error {
	for vi, _ := range jsonSchemaDef.Fields {
		for i, _ := range jsonSchemaDef.Fields[vi].EvalMetrics {
			if jsonSchemaDef.Fields[vi].EvalMetrics[i] == nil {
				jsonSchemaDef.Fields[vi].EvalMetrics[i] = &artemis_orchestrations.EvalMetric{}
			}
			jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult = &artemis_orchestrations.EvalMetricResult{}
			eocr := artemis_orchestrations.EvalMetaDataResult{
				EvalOpCtxStr:               "",
				Operator:                   jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator,
				EvalMetricComparisonValues: jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues,
				FieldValue:                 &jsonSchemaDef.Fields[vi].FieldValue,
			}
			switch jsonSchemaDef.Fields[vi].DataType {
			case "integer", "int":
				if jsonSchemaDef.Fields[vi].IntegerValue == nil {
					return fmt.Errorf("no int value for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil && jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger == nil {
					return fmt.Errorf("no comparison number for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber != nil {
					fv := aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber)
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger = aws.Int(int(fv))
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber = nil
				}
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger != nil {
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetIntEvalComparisonResult(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, *jsonSchemaDef.Fields[vi].IntegerValue, aws.ToInt(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger)))
					eocr.EvalOpCtxStr = fmt.Sprintf("%s %d %s %d", jsonSchemaDef.Fields[vi].FieldName, aws.ToInt(jsonSchemaDef.Fields[vi].IntegerValue), jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToInt(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger))
				} else if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber != nil {
					eocr.EvalOpCtxStr = fmt.Sprintf("%s %d %s %d", jsonSchemaDef.Fields[vi].FieldName, aws.ToInt(jsonSchemaDef.Fields[vi].IntegerValue), jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, int(aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber)))
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetIntEvalComparisonResult(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToInt(jsonSchemaDef.Fields[vi].IntegerValue), int(aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))))
				} else {
					return fmt.Errorf("no comparison number for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
			case "number", "float":
				if jsonSchemaDef.Fields[vi].NumberValue == nil {
					return fmt.Errorf("no number value for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil {
					return fmt.Errorf("no comparison number for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				eocr.EvalOpCtxStr = fmt.Sprintf("%s %f %s %f", jsonSchemaDef.Fields[vi].FieldName, aws.ToFloat64(jsonSchemaDef.Fields[vi].NumberValue), jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetNumericEvalComparisonResult(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToFloat64(jsonSchemaDef.Fields[vi].NumberValue), aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber)))
			case "string":
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString == nil {
					return fmt.Errorf("no comparison string for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].StringValue == nil {
					return fmt.Errorf("no string value for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				eocr.EvalOpCtxStr = fmt.Sprintf("%s %s %s %s", jsonSchemaDef.Fields[vi].FieldName, aws.ToString(jsonSchemaDef.Fields[vi].StringValue), jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToString(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString))
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetStringEvalComparisonResult(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToString(jsonSchemaDef.Fields[vi].StringValue), aws.ToString(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString)))
			case "boolean":
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean == nil {
					return fmt.Errorf("no comparison boolean for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].BooleanValue == nil {
					return fmt.Errorf("no boolean value for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}

				eocr.EvalOpCtxStr = fmt.Sprintf("%s %t %s %t", jsonSchemaDef.Fields[vi].FieldName, *jsonSchemaDef.Fields[vi].BooleanValue, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, *jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean)
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(GetBooleanEvalComparisonResult(*jsonSchemaDef.Fields[vi].BooleanValue, *jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean))
			case "array[integer]":
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil && jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger == nil {
					return fmt.Errorf("no comparison number for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger != nil {
					results, rerr := EvaluateIntArray(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, jsonSchemaDef.Fields[vi].IntegerValueSlice, *jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger)
					if rerr != nil {
						return rerr
					}
					eocr.EvalOpCtxStr = fmt.Sprintf("%s %d %s %d", jsonSchemaDef.Fields[vi].FieldName, jsonSchemaDef.Fields[vi].IntegerValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToInt(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonInteger))
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
				} else {
					results, rerr := EvaluateIntArray(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, jsonSchemaDef.Fields[vi].IntegerValueSlice, int(*jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
					if rerr != nil {
						return rerr
					}
					eocr.EvalOpCtxStr = fmt.Sprintf("%s %v %s %f", jsonSchemaDef.Fields[vi].FieldName, jsonSchemaDef.Fields[vi].IntegerValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
					jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
				}
			case "array[number]":
				if jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber == nil {
					return fmt.Errorf("no comparison number for key '%s'", jsonSchemaDef.Fields[vi].FieldName)
				}
				results, rerr := EvaluateNumericArray(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, jsonSchemaDef.Fields[vi].NumberValueSlice, aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
				if rerr != nil {
					return rerr
				}
				eocr.EvalOpCtxStr = fmt.Sprintf("%s %f %s %f", jsonSchemaDef.Fields[vi].FieldName, jsonSchemaDef.Fields[vi].NumberValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToFloat64(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonNumber))
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			case "array[string]":
				results, rerr := EvaluateStringArray(jsonSchemaDef.Fields[vi].StringValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToString(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString))
				if rerr != nil {
					return rerr
				}
				eocr.EvalOpCtxStr = fmt.Sprintf("%s %s %s %s", jsonSchemaDef.Fields[vi].FieldName, jsonSchemaDef.Fields[vi].StringValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, aws.ToString(jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonString))
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			case "array[boolean]":
				results, rerr := EvaluateBooleanArray(jsonSchemaDef.Fields[vi].BooleanValueSlice, *jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean)
				if rerr != nil {
					return rerr
				}
				eocr.EvalOpCtxStr = fmt.Sprintf("%s %t %s %t", jsonSchemaDef.Fields[vi].FieldName, jsonSchemaDef.Fields[vi].BooleanValueSlice, jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalOperator, *jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricComparisonValues.EvalComparisonBoolean)
				jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalResultOutcomeBool = aws.Bool(Pass(results))
			default:
				return fmt.Errorf("unknown data type '%s'", jsonSchemaDef.Fields[vi].DataType)
			}
			b, err := json.Marshal(eocr)
			if err != nil {
				log.Err(err).Msg("failed to marshal eval op ctx")
				return err
			}
			jsonSchemaDef.Fields[vi].EvalMetrics[i].EvalMetricResult.EvalMetadata = b
		}
	}
	return nil
}

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
