package artemis_orchestrations

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func ConvertToJsonSchema(fd openai.FunctionDefinition) []*JsonSchemaDefinition {
	var schemas []*JsonSchemaDefinition
	jsd, oks := fd.Parameters.(jsonschema.Definition)
	if !oks {
		log.Error().Msg("failed to convert to jsonschema.Definition")
		return schemas
	}
	for name, def := range jsd.Properties {
		schema := &JsonSchemaDefinition{
			SchemaName: name,
			Fields:     []JsonSchemaField{},
		}
		if def.Type == jsonschema.Array {
			schema.IsObjArray = true
			if def.Items == nil {
				continue
			}
			for fieldName, fdef := range def.Items.Properties {
				ft := jsonSchemaDataType(fdef.Type)
				if fdef.Type == jsonschema.Array && fdef.Items != nil {
					ft = fmt.Sprintf("array[%s]", jsonSchemaDataType(fdef.Items.Type))
				}
				schema.Fields = append(schema.Fields, JsonSchemaField{
					FieldName:        fieldName,
					FieldDescription: fdef.Description,
					DataType:         ft,
				})
			}
			schemas = append(schemas, schema)
		} else {
			for fieldName, fdef := range def.Properties {
				ft := jsonSchemaDataType(fdef.Type)
				if fdef.Type == jsonschema.Array && fdef.Items != nil {
					ft = fmt.Sprintf("array[%s]", jsonSchemaDataType(fdef.Items.Type))
				}
				schema.Fields = append(schema.Fields, JsonSchemaField{
					FieldName:        fieldName,
					FieldDescription: fdef.Description,
					DataType:         ft,
				})
			}
			schemas = append(schemas, schema)
		}
	}
	return schemas
}

// Helper function to convert jsonschema.Type to a string representation
func jsonSchemaDataType(t jsonschema.DataType) string {
	switch t {
	case jsonschema.String:
		return "string"
	case jsonschema.Integer:
		return "integer"
	case jsonschema.Number:
		return "number"
	case jsonschema.Boolean:
		return "boolean"
	case jsonschema.Object:
		return "object"
	case jsonschema.Array:
		return "array"
	default:
		return "unknown"
	}
}

// getFieldNames returns a slice of field names from the slice of JsonSchemaField
func getFieldNames(fields []JsonSchemaField) []string {
	var names []string
	for _, field := range fields {
		names = append(names, field.FieldName)
	}
	return names
}

func jsonSchemaType(dataType string) jsonschema.DataType {
	switch dataType {
	case "string":
		return jsonschema.String
	case "number":
		return jsonschema.Number
	case "boolean":
		return jsonschema.Boolean
	case "integer", "int":
		return jsonschema.Integer
	case "array":
		return jsonschema.Array
	case "object":
		return jsonschema.Object
	default:
		return jsonschema.String // default or throw an error based on your requirements
	}
}

func ConvertToFuncDef(schemas []*JsonSchemaDefinition) openai.FunctionDefinition {
	var fnName string
	for _, schema := range schemas {
		fnName += schema.SchemaName
	}

	fd := openai.FunctionDefinition{
		Name:       fnName,
		Parameters: ConvertToFuncDefJsonSchemas(fnName, schemas), // Set the combined schema here
	}
	return fd
}

func ConvertToFuncDefJsonSchemas(fnName string, schemas []*JsonSchemaDefinition) jsonschema.Definition {
	// Initialize the combined properties
	combinedProperties := make(map[string]jsonschema.Definition)
	// Iterate over each schema and create a field for each
	for _, schema := range schemas {
		schemaField := convertDbJsonSchemaFieldsSchema(fnName, schema)
		// If the schema represents an array of objects, adjust the type and items
		if schema.IsObjArray {
			schemaField = jsonschema.Definition{
				Type:     jsonschema.Array,
				Required: []string{schema.SchemaName},
				Items: &jsonschema.Definition{
					Type:       jsonschema.Object,
					Properties: schemaField.Properties,
					Required:   schemaField.Required,
				},
			}
		}
		// Add this schema field to the combined properties
		combinedProperties[schema.SchemaName] = schemaField
	}
	// Create the combined schema object
	combinedSchema := jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: combinedProperties,
	}

	var requiredFields []string

	for k, _ := range combinedSchema.Properties {
		requiredFields = append(requiredFields, k)
	}
	combinedSchema.Required = requiredFields
	return combinedSchema
}

func convertDbJsonSchemaFieldsSchema(fnName string, schema *JsonSchemaDefinition) jsonschema.Definition {
	if schema == nil {
		return jsonschema.Definition{}
	}
	properties := make(map[string]jsonschema.Definition)
	var requiredFields []string
	for _, field := range schema.Fields {
		fieldDef := jsonschema.Definition{
			Description: field.FieldDescription,
		}
		switch field.DataType {
		case "array[number]":
			fieldDef.Type = jsonschema.Array
			fieldDef.Items = &jsonschema.Definition{Type: jsonschema.Number}
		case "array[string]":
			fieldDef.Type = jsonschema.Array
			fieldDef.Items = &jsonschema.Definition{Type: jsonschema.String}
		case "array[boolean]":
			fieldDef.Type = jsonschema.Array
			fieldDef.Items = &jsonschema.Definition{Type: jsonschema.Boolean}
		case "array[integer]":
			fieldDef.Type = jsonschema.Array
			fieldDef.Items = &jsonschema.Definition{Type: jsonschema.Integer}
		default:
			fieldDef.Type = jsonSchemaType(field.DataType) // Assume this function correctly returns the jsonschema type
		}
		properties[field.FieldName] = fieldDef
	}
	for k, _ := range properties {
		requiredFields = append(requiredFields, k)
	}
	return jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   requiredFields,
	}
}

func AssignMapValuesMultipleJsonSchemasSlice(szs []*JsonSchemaDefinition, ms any) ([]JsonSchemaDefinition, error) {
	var responses []JsonSchemaDefinition
	for _, sv := range szs {
		mis, ok := ms.([]map[string]interface{})
		msng, ook := ms.(map[string]interface{})
		if sv == nil {
			continue
		}
		sd := *sv
		if ok {
			for _, inVal := range mis {
				resp, err := AssignMapValuesJsonSchemaFieldsSlice(sd, inVal)
				if err != nil {
					log.Err(err).Interface("inVal", inVal).Msg("AssignMapValuesMultipleJsonSchemasSlice: AssignMapValuesJsonSchemaFieldsSlice failed")
					return nil, err
				}
				if resp != nil {
					responses = append(responses, resp...)
				}
			}
		} else if ook {
			resp, err := AssignMapValuesJsonSchemaFieldsSlice(sd, msng)
			if err != nil {
				log.Err(err).Interface("msng", msng).Msg("AssignMapValuesMultipleJsonSchemasSlice: AssignMapValuesJsonSchemaFieldsSlice failed")
				return nil, err
			}
			if resp != nil {
				responses = append(responses, resp...)
			}
		}
	}

	if !ValidateSchemas(responses) {
		return nil, fmt.Errorf("AssignMapValuesMultipleJsonSchemasSlice: failed to validate schemas")
	}
	return responses, nil
}

func CopyJsonSchemaFieldsSlice(sz JsonSchemaDefinition) JsonSchemaDefinition {
	var tmp []JsonSchemaField
	for _, f := range sz.Fields {
		tmp = append(tmp, JsonSchemaField{
			FieldID:    f.FieldID,
			FieldStrID: f.FieldStrID,
			FieldName:  f.FieldName,
			//FieldDescription: f.FieldDescription, this adds too much extra data
			DataType: f.DataType,
			FieldValue: FieldValue{
				IntegerValue:      nil,
				StringValue:       nil,
				NumberValue:       nil,
				BooleanValue:      nil,
				IntegerValueSlice: nil,
				StringValueSlice:  nil,
				NumberValueSlice:  nil,
				BooleanValueSlice: nil,
				IsValidated:       false,
			},
			EvalMetrics: f.EvalMetrics,
		})
	}

	return JsonSchemaDefinition{
		SchemaID:          sz.SchemaID,
		SchemaStrID:       sz.SchemaStrID,
		SchemaName:        sz.SchemaName,
		SchemaGroup:       sz.SchemaGroup,
		SchemaDescription: sz.SchemaDescription,
		IsObjArray:        sz.IsObjArray,
		Fields:            tmp,
		FieldsMap:         sz.FieldsMap,
	}
}

func AssignMapValuesJsonSchemaFieldsSlice(sz JsonSchemaDefinition, m map[string]interface{}) ([]JsonSchemaDefinition, error) {
	var schemas []JsonSchemaDefinition
	if sz.IsObjArray {
		// Handle case where sz is an array of objects
		for _, inVal := range m {
			vi, vok := inVal.([]interface{})
			if vok {
				for _, inVal2 := range vi {
					vmi, bok := inVal2.(map[string]interface{})
					// Check if the map contains sz.SchemaName as a key
					if bok {
						tmp := CopyJsonSchemaFieldsSlice(sz)
						err := AssignMapValuesJsonSchemaFields(tmp.Fields, vmi)
						if err != nil {
							log.Err(err).Interface("vmi", vmi).Msg("1_AssignMapValuesJsonSchemaFieldsSlice: AssignMapValuesJsonSchemaFields failed")
							return nil, err
						}
						schemas = append(schemas, tmp)

					}
				}
			}
		}
	} else {
		// Handle case where sz is a single object
		// Check if the map contains sz.SchemaName as a key
		if vfi, found := m[sz.SchemaName]; found {
			vfim, vfiok := vfi.(map[string]interface{})
			vi, vok := vfi.([]interface{})
			if vfiok {
				tmp := CopyJsonSchemaFieldsSlice(sz)
				err := AssignMapValuesJsonSchemaFields(tmp.Fields, vfim)
				if err != nil {
					log.Err(err).Interface("vfim", vfim).Msg("2_AssignMapValuesJsonSchemaFieldsSlice: AssignMapValuesJsonSchemaFields failed")
					return nil, err
				}
				schemas = append(schemas, tmp)

			}
			if vok {
				for _, inVal := range vi {
					vmi, bok := inVal.(map[string]interface{})
					if bok {
						tmp := CopyJsonSchemaFieldsSlice(sz)
						err := AssignMapValuesJsonSchemaFields(sz.Fields, vmi)
						if err != nil {
							log.Err(err).Interface("vmi", vmi).Msg("3_AssignMapValuesJsonSchemaFieldsSlice: AssignMapValuesJsonSchemaFields failed")
							return nil, err
						}
						schemas = append(schemas, tmp)
					}
				}
			}
		}
	}
	return schemas, nil
}

func AssignMapValuesJsonSchemaFields(fields []JsonSchemaField, m map[string]interface{}) error {
	if len(m) == 0 {
		return nil
	}
	for i, _ := range fields {
		if val, ok1 := m[fields[i].FieldName]; ok1 {
			switch fields[i].DataType {
			case "string":
				if strVal, okStr := val.(string); okStr {
					fields[i].StringValue = &strVal
					//fmt.Printf("Field %s is a string: %s\n", fieldDef.FieldName, strVal)
					fields[i].IsValidated = true
				} else {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to string", val)
				}
			case "integer":
				if intVal, okInt := val.(int); okInt {
					fields[i].IntegerValue = &intVal
					fields[i].IsValidated = true
					fmt.Printf("Field %s is an integer: %d\n", fields[i].FieldName, intVal)
				} else if fintVal, okfintVal := val.(float64); okfintVal {
					fields[i].IntegerValue = aws.Int(int(fintVal))
					fields[i].IsValidated = true
					//fmt.Printf("Field %s is an float -> integer: %d\n", fieldDef.FieldName, intVal)
				} else {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to int", val)
				}
			case "number":
				if numVal, okNum := val.(float64); okNum {
					fields[i].NumberValue = &numVal
					fields[i].IsValidated = true
					fmt.Printf("Field %s is a number: %f\n", fields[i].FieldName, numVal)
				} else if numValInt, okNumInt := val.(int); okNumInt {
					numValFloat := float64(numValInt)
					fields[i].NumberValue = &numValFloat
					fields[i].IsValidated = true
					//fmt.Printf("Field %s is a number: %f\n", fieldDef.FieldName, numValFloat)
				} else {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to float64", val)
				}
			case "boolean":
				if boolVal, okBool := val.(bool); okBool {
					fields[i].BooleanValue = &boolVal
					fields[i].IsValidated = true
					fmt.Printf("Field %s is a boolean: %t\n", fields[i].FieldName, boolVal)
				} else {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to bool", val)
				}
			case "array[number]":
				vin, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []integer", val)
				}
				vfs, err := interfaceSliceToFloat64Slice(vin)
				if err != nil {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []integer", val)
				}
				fields[i].IsValidated = true
				fields[i].NumberValueSlice = vfs
			case "array[integer]":
				vin, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []integer", val)
				}
				vins, err := interfaceSliceToIntSlice(vin)
				if err != nil {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []integer", val)
				}
				fields[i].IsValidated = true
				fields[i].IntegerValueSlice = vins
			case "array[string]":
				vin, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []string", val)
				}
				vins, err := interfaceSliceToStringSlice(vin)
				if err != nil {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []string", val)
				}
				fields[i].IsValidated = true
				fields[i].StringValueSlice = vins
			case "array[boolean]":
				vin, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []boolean", val)
				}
				bs, err := interfaceSliceToBoolSlice(vin)
				if err != nil {
					return fmt.Errorf("AssignMapValuesJsonSchemaFields: failed to convert %v to []integer", val)
				}
				fields[i].IsValidated = true
				fields[i].BooleanValueSlice = bs
			}
		}
	}
	return nil
}
func interfaceSliceToIntSlice(interfaceSlice []interface{}) ([]int, error) {
	intSlice := make([]int, len(interfaceSlice))
	for i, v := range interfaceSlice {
		// Try asserting as int
		if intValue, ok := v.(int); ok {
			intSlice[i] = intValue
		} else if floatValue, fok := v.(float64); fok {
			// Convert float64 to int
			intSlice[i] = int(floatValue)
		} else {
			return nil, fmt.Errorf("value at index %d is neither int nor float64", i)
		}
	}
	return intSlice, nil
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

func CreateMapInterfaceFromAssignedSchemaFields(schemas []JsonSchemaDefinition) []map[string]interface{} {
	var resp []map[string]interface{}

	for _, schema := range schemas {
		fieldMap := make(map[string]interface{})
		for _, field := range schema.Fields {
			switch field.DataType {
			case "string":
				if field.StringValue != nil {
					fieldMap[field.FieldName] = *field.StringValue
				}
			case "number":
				if field.NumberValue != nil {
					fieldMap[field.FieldName] = *field.NumberValue
				}
			case "integer":
				if field.IntegerValue != nil {
					fieldMap[field.FieldName] = *field.IntegerValue
				}
			case "boolean":
				if field.BooleanValue != nil {
					fieldMap[field.FieldName] = *field.BooleanValue
				}
			case "array[number]":
				if field.NumberValueSlice != nil {
					fieldMap[field.FieldName] = field.NumberValueSlice
				}
			case "array[string]":
				if field.StringValueSlice != nil {
					fieldMap[field.FieldName] = field.StringValueSlice
				}
			case "array[boolean]":
				if field.BooleanValueSlice != nil {
					fieldMap[field.FieldName] = field.BooleanValueSlice
				}
			case "array[integer]":
				if field.IntegerValueSlice != nil {
					fieldMap[field.FieldName] = field.IntegerValueSlice
				}
			}
		}
		resp = append(resp, fieldMap)
	}
	return resp
}
