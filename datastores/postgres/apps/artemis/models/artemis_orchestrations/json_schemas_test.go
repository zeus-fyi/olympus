package artemis_orchestrations

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestConvertToFuncDef3() {
	schema := JsonSchemaDefinition{
		SchemaName: "TestSchema",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "id", DataType: "string", FieldDescription: "id"},
			{FieldName: "values", DataType: "array[number]", FieldDescription: "values"},
		},
	}
	schema2 := JsonSchemaDefinition{
		SchemaName: "Scoring",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "score", DataType: "number", FieldDescription: "scores"},
			{FieldName: "products", DataType: "array[string]", FieldDescription: "products"},
		},
	}
	fd := ConvertToFuncDef("test", []JsonSchemaDefinition{schema, schema2})
	s.Require().NotNil(fd, "Failed to convert JSON schema to OpenAI function definition")
}

func (s *OrchestrationsTestSuite) TestConvertToFuncDef4() {
	schema := JsonSchemaDefinition{
		SchemaName: "lead_scoring",
		IsObjArray: false,
		Fields: []JsonSchemaField{
			{FieldName: "msg_ids", DataType: "array[number]", FieldDescription: "system message ids"},
		},
	}
	fd := ConvertToFuncDef("twitter_extract_tweets", []JsonSchemaDefinition{schema})

	fd2 := FilterAndExtractRelevantTweetsJsonSchemaFunctionDef("system message ids")
	fmt.Println(fd2)
	s.Require().NotNil(fd, "Failed to convert JSON schema to OpenAI function definition")
}

func FilterAndExtractRelevantTweetsJsonSchemaFunctionDef(keepMsgInst string) openai.FunctionDefinition {
	properties := make(map[string]jsonschema.Definition)
	keepMsgs := jsonschema.Definition{
		Type:        jsonschema.Array,
		Description: keepMsgInst,
		Items: &jsonschema.Definition{
			Type: jsonschema.Number,
		},
	}
	properties["msg_ids"] = keepMsgs
	fdSchema := jsonschema.Definition{
		Type:       jsonschema.Object,
		Properties: properties,
		Required:   []string{"msg_ids"},
	}
	fd := openai.FunctionDefinition{
		Name:       "twitter_extract_tweets",
		Parameters: fdSchema,
	}
	return fd
}

func (s *OrchestrationsTestSuite) TestInsertJsonSchema() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID

	// Define a new JSON schema
	newSchema := &JsonSchemaDefinition{
		SchemaName: "testSchema",
		IsObjArray: false,
		Fields: []JsonSchemaField{
			{FieldName: "field1", DataType: "string"},
			{FieldName: "field2", DataType: "int"},
		},
	}

	// Insert the new schema
	err := CreateOrUpdateJsonSchema(ctx, ou, newSchema, nil)
	s.Require().NoError(err, "Failed to insert new JSON schema")

	// Update the schema by removing a field, modifying a field, and adding a new field
	updatedSchema := &JsonSchemaDefinition{
		SchemaName: "testSchema",
		IsObjArray: true, // Example of modifying schema property
		Fields: []JsonSchemaField{
			{FieldName: "field2", DataType: "integer"}, // Modified data type
			{FieldName: "field3", DataType: "boolean"}, // New field
		},
	}
	// Update the existing schema
	err = CreateOrUpdateJsonSchema(ctx, ou, updatedSchema, nil)
	s.Require().NoError(err, "Failed to update JSON schema")
}

func (s *OrchestrationsTestSuite) TestSelectJsonSchemas() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID

	js, err := SelectJsonSchemaByOrg(ctx, ou)
	s.Require().NoError(err, "Failed to select JSON schemas")
	s.Require().NotNil(js, "Failed to select JSON schemas")

}

func (s *OrchestrationsTestSuite) TestJsonParsing() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	js, err := SelectJsonSchemaByOrg(ctx, ou)
	s.Require().NoError(err, "Failed to select JSON schemas")
	v := ConvertToFuncDef("fn", js)
	s.Require().NotNil(v, "Failed to convert JSON schema to OpenAI function definition")
}
