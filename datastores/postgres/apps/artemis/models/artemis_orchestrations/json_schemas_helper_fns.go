package artemis_orchestrations

import (
	"fmt"

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
	// Iterate through the properties of the FunctionDefinition
	for name, def := range jsd.Properties {
		schema := &JsonSchemaDefinition{
			SchemaName: name,
			Fields:     []JsonSchemaField{},
		}
		// Check if the property is an array of objects
		if def.Type == jsonschema.Array && def.Items != nil {
			schema.IsObjArray = true
			for fieldName, fieldDef := range def.Items.Properties {
				schema.Fields = append(schema.Fields, JsonSchemaField{
					FieldName:        fieldName,
					FieldDescription: fieldDef.Description,
					DataType:         jsonSchemaDataType(fieldDef.Type),
				})
			}
		} else if def.Type == jsonschema.Object {
			// Property is an object
			for fieldName, fieldDef := range def.Properties {
				schema.Fields = append(schema.Fields, JsonSchemaField{
					FieldName:        fieldName,
					FieldDescription: fieldDef.Description,
					DataType:         jsonSchemaDataType(fieldDef.Type),
				})
			}
		} else {
			// Simple field
			schema.Fields = append(schema.Fields, JsonSchemaField{
				FieldName:        name,
				FieldDescription: def.Description,
				DataType:         jsonSchemaDataType(def.Type),
			})
		}
		schemas = append(schemas, schema)
	}
	return schemas
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
	case "array":
		return jsonschema.Array
	case "object":
		return jsonschema.Object
	default:
		return jsonschema.String // default or throw an error based on your requirements
	}
}

const (
	msgID                                = "msg_id"
	analyzedMsgId                        = "analyzed_msg_id"
	socialMediaEngagementResponseFormat  = "social-media-engagement"
	keepTweetRelationshipToSingleMessage = "add the msg_id from the msg_body field that you are analyzing"
)

func ConvertToFuncDef(fnName string, schemas []*JsonSchemaDefinition) openai.FunctionDefinition {
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
		schema.IsObjArray = true
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

	if fnName == socialMediaEngagementResponseFormat {
		properties[analyzedMsgId] = jsonschema.Definition{
			Type:        jsonschema.Number,
			Description: keepTweetRelationshipToSingleMessage,
		}
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

func AssignMapValuesMultipleJsonSchemasSlice(szs []*JsonSchemaDefinition, ms any) [][]*JsonSchemaDefinition {
	var responses [][]*JsonSchemaDefinition

	for _, sz := range szs {
		mis, ok := ms.([]map[string]interface{})
		msng, ook := ms.(map[string]interface{})
		if ok {
			for _, mi := range mis {
				responses = append(responses, AssignMapValuesJsonSchemaFieldsSlice(sz, mi))
			}
		} else if ook {
			responses = append(responses, AssignMapValuesJsonSchemaFieldsSlice(sz, msng))
		}
	}
	return responses
}

func AssignMapValuesJsonSchemaFieldsSlice(sz *JsonSchemaDefinition, m map[string]interface{}) []*JsonSchemaDefinition {
	if sz == nil {
		return nil
	}
	var schemas []*JsonSchemaDefinition
	if sz.IsObjArray {
		// Handle case where sz is an array of objects
		for _, v := range m {
			vi, vok := v.([]interface{})
			if vok {
				for _, item := range vi {
					vmi, bok := item.(map[string]interface{})
					// Check if the map contains sz.SchemaName as a key
					if bok {
						jsd := AssignMapValuesJsonSchemaFields(sz, vmi)
						schemas = append(schemas, jsd)
					}
				}
			}
		}
	} else {
		// Handle case where sz is a single object
		// Check if the map contains sz.SchemaName as a key
		if vfi, found := m[sz.SchemaName]; found {
			vfim, vfiok := vfi.(map[string]interface{})
			if vfiok {
				jsd := AssignMapValuesJsonSchemaFields(sz, vfim)
				schemas = append(schemas, jsd)
			}
		}
	}
	return schemas
}

//func AssignMapValuesMultipleJsonSchemasSlice(szs []*JsonSchemaDefinition, ms any) [][]*JsonSchemaDefinition {
//	var responses [][]*JsonSchemaDefinition
//
//	for _, sz := range szs {
//		mis, ok := ms.([]map[string]interface{})
//		msng, ook := ms.(map[string]interface{})
//		if ok {
//			for _, mi := range mis {
//				responses = append(responses, AssignMapValuesJsonSchemaFieldsSlice(sz, mi))
//			}
//		} else if ook {
//			responses = append(responses, AssignMapValuesJsonSchemaFieldsSlice(sz, msng))
//		}
//	}
//	return responses
//}
//
//func AssignMapValuesJsonSchemaFieldsSlice(sz *JsonSchemaDefinition, m any) []*JsonSchemaDefinition {
//	if sz == nil {
//		return nil
//	}
//	var schemas []*JsonSchemaDefinition
//	if sz.IsObjArray {
//		// Handle case where sz is an array of objects
//		sliceOfMaps, ok := m.(map[string]interface{})
//		if !ok {
//			return nil // or handle the error as you see fit
//		}
//		for _, v := range sliceOfMaps {
//			vi, vok := v.([]interface{})
//			if vok {
//				for i, _ := range vi {
//					vmi, bok := vi[i].(map[string]interface{})
//					if bok {
//						jsd := AssignMapValuesJsonSchemaFields(sz, vmi)
//						schemas = append(schemas, jsd)
//					}
//				}
//			}
//		}
//	} else {
//		// Handle case where sz is a single object
//		jsd := AssignMapValuesJsonSchemaFields(sz, m.(map[string]interface{}))
//		schemas = append(schemas, jsd)
//	}
//	return schemas
//}

func AssignMapValuesJsonSchemaFields(sz *JsonSchemaDefinition, m map[string]interface{}) *JsonSchemaDefinition {
	if sz == nil || len(m) == 0 {
		return nil
	}
	for i, _ := range sz.Fields {
		fieldDef := &sz.Fields[i] // Get a reference to the field definition
		if val, ok1 := m[fieldDef.FieldName]; ok1 {
			switch fieldDef.DataType {
			case "string":
				if strVal, okStr := val.(string); okStr {
					fieldDef.StringValue = &strVal
					fmt.Printf("Field %s is a string: %s\n", fieldDef.FieldName, strVal)
				}
			case "integer":
				if intVal, okInt := val.(int); okInt {
					fieldDef.IntValue = &intVal
					fmt.Printf("Field %s is an integer: %d\n", fieldDef.FieldName, intVal)
				}
			case "number":
				if numVal, okNum := val.(float64); okNum {
					fieldDef.NumberValue = &numVal
					fmt.Printf("Field %s is a number: %f\n", fieldDef.FieldName, numVal)
				} else if numValInt, okNumInt := val.(int); okNumInt {
					numValFloat := float64(numValInt)
					fieldDef.NumberValue = &numValFloat
					fmt.Printf("Field %s is a number: %f\n", fieldDef.FieldName, numValFloat)
				}
			case "boolean":
				if boolVal, okBool := val.(bool); okBool {
					fieldDef.BooleanValue = &boolVal
					fmt.Printf("Field %s is a boolean: %t\n", fieldDef.FieldName, boolVal)
				}
			case "array[number]":
				if sliceVal, okArrayNum := val.([]interface{}); okArrayNum {
					numbers := make([]float64, 0)
					for _, v := range sliceVal {
						if num, okNum := v.(float64); okNum {
							numbers = append(numbers, num)
						}
					}
					fieldDef.NumberValueSlice = numbers
					fmt.Printf("Field %s is an array of numbers: %v\n", fieldDef.FieldName, numbers)
				}
			case "array[integer]":
				if sliceVal, okArrayInt := val.([]interface{}); okArrayInt {
					intSlice := make([]int, 0)
					for _, v := range sliceVal {
						if str, okInt := v.(int); okInt {
							intSlice = append(intSlice, str)
						}
					}
					fieldDef.IntValueSlice = intSlice
					fmt.Printf("Field %s is an array of ints: %v\n", fieldDef.FieldName, intSlice)
				}
			case "array[string]":
				if sliceVal, okArrayStr := val.([]interface{}); okArrayStr {
					strings := make([]string, 0)
					for _, v := range sliceVal {
						if str, okStr := v.(string); okStr {
							strings = append(strings, str)
						}
					}
					fieldDef.StringValueSlice = strings
					fmt.Printf("Field %s is an array of strings: %v\n", fieldDef.FieldName, strings)
				}
			case "array[boolean]":
				if sliceVal, okArrayBool := val.([]interface{}); okArrayBool {
					booleans := make([]bool, 0)
					for _, v := range sliceVal {
						if b, okBool := v.(bool); okBool {
							booleans = append(booleans, b)
						}
					}
					fieldDef.BooleanValueSlice = booleans
					fmt.Printf("Field %s is an array of booleans: %v\n", fieldDef.FieldName, booleans)
				}
			}
		}
	}
	return sz
}
