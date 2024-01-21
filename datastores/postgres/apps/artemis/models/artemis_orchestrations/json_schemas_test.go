package artemis_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestConvertToFuncDef3() {
	schema := &JsonSchemaDefinition{
		SchemaName: "messages",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "id", DataType: "string", FieldDescription: "id"},
			{FieldName: "values", DataType: "array[number]", FieldDescription: "values"},
		},
	}
	schema2 := &JsonSchemaDefinition{
		SchemaName: "scoring",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "score", DataType: "number", FieldDescription: "scores"},
			{FieldName: "products", DataType: "array[string]", FieldDescription: "products"},
		},
	}
	fd := ConvertToFuncDef("test", []*JsonSchemaDefinition{schema, schema2})
	s.Require().NotNil(fd, "Failed to convert JSON schema to OpenAI function definition")
}

func (s *OrchestrationsTestSuite) TestConvertToFuncDef4() {
	schema := &JsonSchemaDefinition{
		SchemaName: "lead_scoring",
		IsObjArray: false,
		Fields: []JsonSchemaField{
			{FieldName: "msg_ids", DataType: "array[number]", FieldDescription: "system message ids"},
		},
	}
	fd := ConvertToFuncDef("twitter_extract_tweets", []*JsonSchemaDefinition{schema})

	s.Require().NotNil(fd, "Failed to convert JSON schema to OpenAI function definition")
}
func (s *OrchestrationsTestSuite) TestInsertJsonSchema2() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID

	schema := JsonSchemaDefinition{
		SchemaName: "web3-sales-lead-scoring-2",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "aggregate_lead_score", DataType: "number", FieldDescription: "aggregate_lead_score description"},
			{FieldName: "lead_score_metrics", DataType: "array[string]", FieldDescription: "lead_score_metrics description"},
		},
	}

	// Insert the new schema
	err := CreateOrUpdateJsonSchema(ctx, ou, &schema, nil)
	s.Require().NoError(err, "Failed to insert new JSON schema")
}
func (s *OrchestrationsTestSuite) TestInsertJsonSchema() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID

	schema := JsonSchemaDefinition{
		SchemaName: "messages",
		IsObjArray: true,
		Fields: []JsonSchemaField{
			{FieldName: "id", DataType: "string", FieldDescription: "id"},
			{FieldName: "values", DataType: "array[number]", FieldDescription: "values"},
		},
	}

	// Insert the new schema
	err := CreateOrUpdateJsonSchema(ctx, ou, &schema, nil)
	s.Require().NoError(err, "Failed to insert new JSON schema")

	// Update the schema by removing a field, modifying a field, and adding a new field
	schema = JsonSchemaDefinition{
		SchemaName: "messages",
		IsObjArray: false,
		Fields: []JsonSchemaField{
			{FieldName: "msg_id", DataType: "number", FieldDescription: "id"},
		},
	}
	// Update the existing schema
	err = CreateOrUpdateJsonSchema(ctx, ou, &schema, nil)
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
	v := ConvertToFuncDef("fn", js.Slice)
	s.Require().NotNil(v, "Failed to convert JSON schema to OpenAI function definition")

	jsb := ConvertToJsonSchema(v)
	s.Require().NotNil(jsb, "Failed to convert OpenAI function definition to JSON schema")

	s.Require().NoError(err, "Failed to convert JSON schema to map")
	//
	//mi := CreateMapInterfaceJson(jsb)
	//s.Require().NotNil(mi, "Failed to convert JSON schema to map interface")
	//fmt.Println(mi)
}
func (s *OrchestrationsTestSuite) TestAssignSliceMapValuesJsonSchemaFields() {
	// This is the map structure based on the provided image
	m := map[string]interface{}{
		"web3-sales-lead-scoring-2": []interface{}{
			map[string]interface{}{
				"analyzed_msg_id":      1704249180,
				"aggregate_lead_score": 15,
				"lead_score_metrics":   "Kubernetes mention, indicating direct interest in technology Zeusfyi enhances.",
			},
			map[string]interface{}{
				"analyzed_msg_id":      2704249180,
				"aggregate_lead_score": 10,
				"lead_score_metrics":   "Other mention, indicating indirect interest in technology Zeusfyi enhances.",
			},
			// Add other map elements as needed based on the actual data
		},
		"add-on": []interface{}{
			map[string]interface{}{
				"test_id":    "test_id_1",
				"test_score": "100",
			},
		},
	}
	sz2 := JsonSchemaDefinition{
		SchemaID:    0,
		SchemaName:  "add-on",
		SchemaGroup: "",
		IsObjArray:  true, // This should be true to match the array structure in the image
		Fields: []JsonSchemaField{
			// Define your schema fields based on the actual structure of the map elements
			{FieldName: "test_id", DataType: "string", FieldDescription: "Analyzed test ID"},
			{FieldName: "test_id", DataType: "string", FieldDescription: "Analyzed test score ID"},
			// Add other fields as needed
		},
	}

	sz := JsonSchemaDefinition{
		SchemaID:    0,
		SchemaName:  "web3-sales-lead-scoring-2",
		SchemaGroup: "",
		IsObjArray:  true, // This should be true to match the array structure in the image
		Fields: []JsonSchemaField{
			// Define your schema fields based on the actual structure of the map elements
			{FieldName: "analyzed_msg_id", DataType: "number", FieldDescription: "Analyzed message ID"},
			{FieldName: "aggregate_lead_score", DataType: "integer", FieldDescription: "Aggregate lead score"},
			{FieldName: "lead_score_metrics", DataType: "string", FieldDescription: "Lead score metrics"},
			// Add other fields as needed
		},
	}

	// Pass the value part of the map to the function, not the entire map
	jr := AssignMapValuesJsonSchemaFieldsSlice(&sz, m["web3-sales-lead-scoring-2"])
	for _, r := range jr {
		fmt.Println(*r)
	}

	// Pass the value part of the map to the function, not the entire map
	jr2 := AssignMapValuesJsonSchemaFieldsSlice(&sz2, m["add-on"])
	for _, r := range jr2 {
		fmt.Println(*r)
	}

	//jr3 := AssignMapValuesMultipleJsonSchemasSlice([]*JsonSchemaDefinition{&sz, &sz2}, m)
	//for _, r := range jr3 {
	//	for _, r2 := range r {
	//		fmt.Println(*r2)
	//	}
	//}
}

func (s *OrchestrationsTestSuite) TestAssignMapValuesJsonSchemaFields() {
	m := map[string]interface{}{
		"id":     "stringID",
		"values": []interface{}{1, 2},
	}
	sz := JsonSchemaDefinition{
		SchemaID:    0,
		SchemaName:  "schema_field",
		SchemaGroup: "",
		IsObjArray:  false,
		Fields: []JsonSchemaField{
			{FieldName: "id", DataType: "string", FieldDescription: "id"},
			{FieldName: "values", DataType: "array[number]", FieldDescription: "values"},
			// Add other field types as needed
			{FieldName: "integer_field", DataType: "integer", FieldDescription: "integer value"},
			{FieldName: "number_field", DataType: "number", FieldDescription: "number value"},
		},
	}
	jr := AssignMapValuesJsonSchemaFields(&sz, m)
	fmt.Println(*jr)
}
