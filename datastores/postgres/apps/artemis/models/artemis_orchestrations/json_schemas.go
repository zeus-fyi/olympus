package artemis_orchestrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type JsonSchemaDefinition struct {
	SchemaID    int               `db:"schema_id" json:"schemaID"`
	SchemaName  string            `db:"schema_name" json:"schemaName"`
	SchemaGroup string            `db:"schema_group" json:"schemaGroup"`
	IsObjArray  bool              `db:"is_obj_array" json:"isObjArray"`
	Fields      []JsonSchemaField `db:"-" json:"fields"`
}

type JsonSchemaField struct {
	FieldID           int         `db:"field_id" json:"fieldID"`
	FieldName         string      `db:"field_name" json:"fieldName"`
	FieldDescription  string      `db:"field_description" json:"fieldDescription"`
	DataType          string      `db:"data_type" json:"dataType"`
	IntValue          *int        `db:"-" json:"intValue,omitempty"`
	StringValue       *string     `db:"-" json:"stringValue,omitempty"`
	NumberValue       *float64    `db:"-" json:"numberValue,omitempty"`
	BooleanValue      *bool       `db:"-" json:"booleanValue,omitempty"`
	IntValueSlice     []int       `db:"-" json:"intValueSlice,omitempty"`
	StringValueSlice  []string    `db:"-" json:"stringValueSlice,omitempty"`
	NumberValueSlice  []float64   `db:"-" json:"numberValueSlice,omitempty"`
	BooleanValueSlice []bool      `db:"-" json:"booleanValueSlice,omitempty"`
	EvalMetric        *EvalMetric `db:"-" json:"evalMetric,omitempty"`
}

func AssignMapValuesMultipleJsonSchemasSlice(szs []*JsonSchemaDefinition, ms any) [][]*JsonSchemaDefinition {
	var responses [][]*JsonSchemaDefinition
	mis, ok := ms.([]map[string]interface{})
	msng, ook := ms.(map[string]interface{})
	for _, sz := range szs {
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

func AssignMapValuesJsonSchemaFieldsSlice(sz *JsonSchemaDefinition, m any) []*JsonSchemaDefinition {
	if sz == nil {
		return nil
	}
	var schemas []*JsonSchemaDefinition
	if sz.IsObjArray {
		// Handle case where sz is an array of objects
		sliceOfMaps, ok := m.(map[string]interface{})
		if !ok {
			return nil // or handle the error as you see fit
		}
		for _, v := range sliceOfMaps {
			vi, vok := v.([]interface{})
			if vok {
				for i, _ := range vi {
					vmi, bok := vi[i].(map[string]interface{})
					if bok {
						jsd := AssignMapValuesJsonSchemaFields(sz, vmi)
						schemas = append(schemas, jsd)
					}
				}
			}
		}
	} else {
		// Handle case where sz is a single object
		jsd := AssignMapValuesJsonSchemaFields(sz, m.(map[string]interface{}))
		schemas = append(schemas, jsd)
	}
	return schemas
}

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
			case "array[int]":
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

type AITaskJsonSchema struct {
	SchemaID int `db:"schema_id" json:"schemaID"`
	TaskID   int `db:"task_id" json:"taskID"`
}

func CreateOrUpdateJsonSchema(ctx context.Context, ou org_users.OrgUser, schema *JsonSchemaDefinition, taskID *int) error {
	if schema == nil {
		return nil
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("failed to start transaction")
		return err
	}
	defer tx.Rollback(ctx)
	var schemaID int

	// Step 1: Check for existing schema in ai_json_schema_definitions
	err = tx.QueryRow(ctx, `
		SELECT schema_id FROM public.ai_json_schema_definitions
		WHERE org_id = $1 AND schema_name = $2;
	`, ou.OrgID, schema.SchemaName).Scan(&schemaID)

	if err == pgx.ErrNoRows {
		// No existing schema, need to create a new one
		err = tx.QueryRow(ctx, `
        INSERT INTO public.ai_schemas (org_id)
        VALUES ($1)
        RETURNING schema_id;
    `, ou.OrgID).Scan(&schemaID)
		if err != nil {
			log.Err(err).Msg("failed to insert into ai_schemas")
			return err
		}
	} else if err != nil {
		log.Err(err).Msg("failed to query existing schema")
		return err
	}

	// Step 2: Insert or Update in ai_json_schema_definitions
	err = tx.QueryRow(ctx, `
    INSERT INTO public.ai_json_schema_definitions(org_id, schema_id, schema_name, is_obj_array)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (org_id, schema_name) DO UPDATE 
    SET is_obj_array = EXCLUDED.is_obj_array, schema_group = EXCLUDED.schema_group
    RETURNING schema_id;
`, ou.OrgID, schemaID, schema.SchemaName, schema.IsObjArray).Scan(&schema.SchemaID)
	if err != nil {
		log.Err(err).Msg("failed to insert or update in json_schema_definitions")
		return err
	}
	// Create a map of field names for easier checking
	fieldMap := make(map[string]bool)
	for _, field := range schema.Fields {
		fieldMap[field.FieldName] = true
	}
	// Archive fields that are not in the new schema
	_, err = tx.Exec(ctx, `
		UPDATE public.ai_fields
		SET is_field_archived = true, archived_at = NOW()
		WHERE schema_id = $1 AND field_name NOT IN (SELECT unnest($2::text[])) AND is_field_archived = false;
	`, schema.SchemaID, pq.Array(getFieldNames(schema.Fields)))
	if err != nil {
		log.Err(err).Msg("failed to archive old fields from ai_fields")
		return err
	}
	// Insert into ai_json_schema_fields
	for _, field := range schema.Fields {
		cr := chronos.Chronos{}

		// Check if the data_type has changed and archive the old field if necessary
		_, err = tx.Exec(ctx, `
        UPDATE public.ai_fields
        SET is_field_archived = true, archived_at = NOW()
        WHERE schema_id = $1 AND field_name = $2 AND data_type != $3 AND is_field_archived = false;
    `, schema.SchemaID, field.FieldName, field.DataType)
		if err != nil {
			log.Err(err).Msg("failed to archive old field in ai_fields")
			return err
		}

		// Insert a new field or update an existing one (if it's not archived)
		_, err = tx.Exec(ctx, `
        INSERT INTO public.ai_fields(field_id, schema_id, field_name, data_type, field_description)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (schema_id, field_name) WHERE is_field_archived = false DO UPDATE 
        SET data_type = EXCLUDED.data_type,
            field_description = EXCLUDED.field_description;
    `, cr.UnixTimeStampNow(), schema.SchemaID, field.FieldName, field.DataType, field.FieldDescription)
		if err != nil {
			log.Err(err).Msg("failed to insert or update field in ai_fields")
			return err
		}
	}

	// Additional step to handle ai_task_json_schemas
	if taskID != nil {
		_, err = tx.Exec(ctx, `
			INSERT INTO public.ai_json_task_schemas(schema_id, task_id)
			VALUES ($1, $2)
			ON CONFLICT (schema_id, task_id) DO NOTHING;
		`, schema.SchemaID, *taskID)
		if err != nil {
			log.Err(err).Msg("failed to insert into ai_task_json_schemas")
			return err
		}
	}
	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
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
		schemaField := convertDbJsonSchemaFieldsTSchema(fnName, schema)
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

func convertDbJsonSchemaFieldsTSchema(fnName string, schema *JsonSchemaDefinition) jsonschema.Definition {
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

type JsonSchemasGroup struct {
	Slice []*JsonSchemaDefinition       `json:"jsonSchemaDefinitionsSlice,omitempty"`
	Map   map[int]*JsonSchemaDefinition `json:"jsonSchemaDefinitionsMap,omitempty"`
}

func SelectJsonSchemaByOrg(ctx context.Context, ou org_users.OrgUser) (*JsonSchemasGroup, error) {
	var schemas []*JsonSchemaDefinition
	// Query to join json_schema_definitions and ai_task_json_schema_fields
	query := `
        SELECT d.schema_id, d.schema_name, d.schema_group, d.is_obj_array, f.field_name, f.data_type, f.field_description
        FROM public.ai_json_schema_definitions d
        JOIN public.ai_json_schema_fields f ON d.schema_id = f.schema_id
        WHERE d.org_id = $1
        ORDER BY d.schema_id, f.field_name;`

	rows, err := apps.Pg.Query(ctx, query, ou.OrgID)
	if err != nil {
		log.Err(err).Msg("failed to execute query for JSON schema")
		return nil, err
	}
	defer rows.Close()

	// A map to keep track of schemas and their fields
	schemaMap := make(map[int]*JsonSchemaDefinition)
	for rows.Next() {
		var field JsonSchemaField
		var schema JsonSchemaDefinition

		err = rows.Scan(&schema.SchemaID, &schema.SchemaName, &schema.SchemaGroup, &schema.IsObjArray, &field.FieldName, &field.DataType, &field.FieldDescription)
		if err != nil {
			log.Err(err).Msg("failed to scan JSON schema row")
			return nil, err
		}

		if s, exists := schemaMap[schema.SchemaID]; exists {
			// If schema already exists in map, append the field to it
			s.Fields = append(s.Fields, field)
		} else {
			// If new schema, initialize and add to map
			schema.Fields = append(schema.Fields, field)
			schemaMap[schema.SchemaID] = &schema
		}
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("failed to iterate over JSON schema rows")
		return nil, err
	}
	// Convert the map to a slice
	for _, schema := range schemaMap {
		schemas = append(schemas, schema)
	}
	return &JsonSchemasGroup{
		Slice: schemas,
		Map:   schemaMap,
	}, nil
}
