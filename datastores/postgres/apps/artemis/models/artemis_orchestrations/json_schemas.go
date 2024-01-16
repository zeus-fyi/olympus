package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type JsonSchemaDefinition struct {
	SchemaID    int               `db:"schema_id" json:"schemaID"`
	SchemaName  string            `db:"schema_name" json:"schemaName"`
	SchemaGroup string            `db:"schema_group" json:"schemaGroup"`
	IsObjArray  bool              `db:"is_obj_array" json:"isObjArray"`
	Fields      []JsonSchemaField `db:"-" json:"fields"`
}

type JsonSchemaField struct {
	FieldName        string `db:"field_name" json:"fieldName"`
	FieldDescription string `db:"field_description" json:"fieldDescription"`
	DataType         string `db:"data_type" json:"dataType"`
}

type AITaskJsonSchema struct {
	SchemaID int `db:"schema_id" json:"schemaID"`
	TaskID   int `db:"task_id" json:"taskID"`
}

func CreateOrUpdateJsonSchema(ctx context.Context, ou org_users.OrgUser, schema *JsonSchemaDefinition, taskID *int) error {
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("failed to start transaction")
		return err
	}
	defer tx.Rollback(ctx)
	// Insert into json_schema_definitions and get the generated schema_id
	err = tx.QueryRow(ctx, `
			INSERT INTO public.ai_json_schema_definitions(org_id, schema_name, is_obj_array)
			VALUES ($1, $2, $3)
			ON CONFLICT (org_id, schema_name) DO UPDATE 
			SET is_obj_array = EXCLUDED.is_obj_array, schema_group = EXCLUDED.schema_group
			RETURNING schema_id;
	`, ou.OrgID, schema.SchemaName, schema.IsObjArray).Scan(&schema.SchemaID)
	if err != nil {
		log.Err(err).Msg("failed to insert or update in json_schema_definitions")
		return err
	}
	// Create a map of field names for easier checking
	fieldMap := make(map[string]bool)
	for _, field := range schema.Fields {
		fieldMap[field.FieldName] = true
	}
	// Delete fields that are not in the new schema
	_, err = tx.Exec(ctx, `
        DELETE FROM public.ai_json_schema_fields 
        WHERE schema_id = $1 AND field_name NOT IN (SELECT unnest($2::text[]));
    `, schema.SchemaID, pq.Array(getFieldNames(schema.Fields)))
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg("failed to delete old fields from ai_task_json_schema_fields")
		return err
	}
	// Insert into ai_json_schema_fields
	for _, field := range schema.Fields {
		_, err = tx.Exec(ctx, `
        INSERT INTO public.ai_json_schema_fields(schema_id, field_name, data_type, field_description)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (schema_id, field_name) DO UPDATE 
        SET data_type = EXCLUDED.data_type, field_description = EXCLUDED.field_description;
    `, schema.SchemaID, field.FieldName, field.DataType, field.FieldDescription)
		if err != nil {
			log.Err(err).Msg("failed to insert or update ai_json_schema_fields")
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

func SelectJsonSchemaByOrg(ctx context.Context, ou org_users.OrgUser) ([]JsonSchemaDefinition, error) {
	var schemas []JsonSchemaDefinition
	// Query to join json_schema_definitions and ai_task_json_schema_fields
	query := `
        SELECT d.schema_id, d.org_id, d.schema_name, d.schema_group, d.is_obj_array, f.field_name, f.data_type, f.field_description
        FROM public.ai_json_schema_definitions d
        JOIN public.ai_task_json_schema_fields f ON d.schema_id = f.schema_id
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
		var schemaID int
		var field JsonSchemaField
		var schema JsonSchemaDefinition

		err = rows.Scan(&schemaID, &ou.OrgID, &schema.SchemaName, &schema.SchemaGroup, &schema.IsObjArray, &field.FieldName, &field.DataType, &field.FieldDescription)
		if err != nil {
			log.Err(err).Msg("failed to scan JSON schema row")
			return nil, err
		}

		if s, exists := schemaMap[schemaID]; exists {
			// If schema already exists in map, append the field to it
			s.Fields = append(s.Fields, field)
		} else {
			// If new schema, initialize and add to map
			schema.Fields = append(schema.Fields, field)
			schemaMap[schemaID] = &schema
		}
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("failed to iterate over JSON schema rows")
		return nil, err
	}
	// Convert the map to a slice
	for _, schema := range schemaMap {
		schemas = append(schemas, *schema)
	}

	return schemas, nil
}
