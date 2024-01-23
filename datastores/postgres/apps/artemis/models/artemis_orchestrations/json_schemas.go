package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
	IntegerValue      *int        `db:"-" json:"intValue,omitempty"`
	StringValue       *string     `db:"-" json:"stringValue,omitempty"`
	NumberValue       *float64    `db:"-" json:"numberValue,omitempty"`
	BooleanValue      *bool       `db:"-" json:"booleanValue,omitempty"`
	IntegerValueSlice []int       `db:"-" json:"intValueSlice,omitempty"`
	StringValueSlice  []string    `db:"-" json:"stringValueSlice,omitempty"`
	NumberValueSlice  []float64   `db:"-" json:"numberValueSlice,omitempty"`
	BooleanValueSlice []bool      `db:"-" json:"booleanValueSlice,omitempty"`
	IsValidated       bool        `db:"-" json:"isValidated,omitempty"`
	EvalMetric        *EvalMetric `db:"-" json:"evalMetricResult,omitempty"`
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
			INSERT INTO public.ai_task_schemas(schema_id, task_id)
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

type JsonSchemasGroup struct {
	Slice []*JsonSchemaDefinition       `json:"jsonSchemaDefinitionsSlice,omitempty"`
	Map   map[int]*JsonSchemaDefinition `json:"jsonSchemaDefinitionsMap,omitempty"`
}

func SelectJsonSchemaByOrg(ctx context.Context, ou org_users.OrgUser) (*JsonSchemasGroup, error) {
	var schemas []*JsonSchemaDefinition
	// Query to join json_schema_definitions and ai_task_json_schema_fields
	query := `
        SELECT d.schema_id, d.schema_name, d.schema_group, d.is_obj_array, f.field_name, f.data_type, f.field_description, f.field_id
        FROM public.ai_json_schema_definitions d
        JOIN public.ai_fields f ON d.schema_id = f.schema_id
        WHERE d.org_id = $1 AND f.is_field_archived = false
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

		err = rows.Scan(&schema.SchemaID, &schema.SchemaName, &schema.SchemaGroup, &schema.IsObjArray, &field.FieldName, &field.DataType, &field.FieldDescription, &field.FieldID)
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
