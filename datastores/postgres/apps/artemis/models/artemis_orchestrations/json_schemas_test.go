package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

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
	v := ConvertToFuncDef("test", js)
	s.Require().NotNil(v, "Failed to convert JSON schema to OpenAI function definition")
}
