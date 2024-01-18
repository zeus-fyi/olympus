package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

func (z *ZeusAiPlatformActivities) AnalyzeEngagementTweets(ctx context.Context, ou org_users.OrgUser, sg *hera_search.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	if sg == nil || sg.SearchResults == nil || len(sg.SearchResults) == 0 {
		log.Info().Msg("AnalyzeEngagementTweets: no search results to analyze engagement")
		return nil, nil
	}
	za := NewZeusAiPlatformActivities()
	params := hera_openai.OpenAIParams{
		Model:              sg.Model,
		FunctionDefinition: sg.FunctionDefinition,
		Prompt:             hera_search.FormatSearchResultsV3(sg.SearchResults),
		SystemPromptExt:    "Analyze the messages using the criteria provided by the schema field definitions.",
	}
	resp, err := za.CreateJsonOutputModelResponse(ctx, ou, params)
	if err != nil {
		log.Err(err).Msg("ExtractTweets: CreateJsonOutputModelResponse failed")
		return nil, err
	}
	m, err := UnmarshallOpenAiJsonInterface(sg.FunctionDefinition.Name, resp)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
		return nil, err
	}
	jsd := artemis_orchestrations.ConvertToJsonSchema(sg.FunctionDefinition)
	err = mapToSchema(m, jsd)
	if err != nil {
		log.Err(err).Msg("AnalyzeEngagementTweets: mapToSchema failed")
		return nil, err
	}
	resp.JsonResponseResults = jsd
	return resp, nil
}

func mapToSchema(m map[string]interface{}, schemas []artemis_orchestrations.JsonSchemaDefinition) error {
	for _, schema := range schemas {
		if schemaValue, ok1 := m[schema.SchemaName]; ok1 {
			switch v := schemaValue.(type) {
			case []interface{}:
				if schema.IsObjArray {
					// Handle array of objects
					for _, item := range v {
						if obj, ok := item.(map[string]interface{}); ok {
							if err := mapObjectToSchemaFields(obj, schema.Fields); err != nil {
								return err
							}
						}
					}
				}
			case map[string]interface{}:
				if !schema.IsObjArray {
					// Handle a single object
					if err := mapObjectToSchemaFields(v, schema.Fields); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func mapObjectToSchemaFields(obj map[string]interface{}, fields []artemis_orchestrations.JsonSchemaField) error {
	for i, field := range fields {
		if fieldValue, ok1 := obj[field.FieldName]; ok1 {
			switch field.DataType {
			case "string":
				if strVal, ok := fieldValue.(string); ok {
					fields[i].StringValue = &strVal
				}
			case "number":
				if numVal, ok := fieldValue.(float64); ok { // JSON numbers are float64
					fields[i].NumberValue = &numVal
				}
			case "boolean":
				if boolVal, ok := fieldValue.(bool); ok {
					fields[i].BooleanValue = &boolVal
				}
			case "array":
				// Assuming arrays of basic types (string, number, boolean)
				switch {
				case field.StringValueSlice != nil:
					if strSlice, ok2 := fieldValue.([]interface{}); ok2 {
						for _, v := range strSlice {
							if str, ok3 := v.(string); ok3 {
								fields[i].StringValueSlice = append(fields[i].StringValueSlice, &str)
							}
						}
					}
				case field.NumberValueSlice != nil:
					if numSlice, ok4 := fieldValue.([]interface{}); ok4 {
						for _, v := range numSlice {
							if num, ok := v.(float64); ok {
								fields[i].NumberValueSlice = append(fields[i].NumberValueSlice, &num)
							}
						}
					}
				case field.BooleanValueSlice != nil:
					if boolSlice, ok5 := fieldValue.([]interface{}); ok5 {
						for _, v := range boolSlice {
							if b, ok := v.(bool); ok {
								fields[i].BooleanValueSlice = append(fields[i].BooleanValueSlice, &b)
							}
						}
					}
				}
			case "object":
				// Assuming the object is another JsonSchemaField or similar structure
				if nestedObj, ok := fieldValue.(map[string]interface{}); ok {
					// You need to define a way to get the fields for the nested object
					// For simplicity, let's assume it's the same fields (recursive structure)
					if err := mapObjectToSchemaFields(nestedObj, fields); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("unsupported data type: %s", field.DataType)
			}
		}
	}
	return nil
}

func UnmarshallFilteredMsgIdsFromAiJsonSmExtraction(fn string, cr *ChatCompletionQueryResponse) error {
	m, err := UnmarshallOpenAiJsonInterface(fn, cr)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
		return err
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: json.Marshal failed")
		return err
	}
	// Unmarshal the JSON string into the FilteredMessages struct
	cr.FilteredMessages = &FilteredMessages{}
	err = json.Unmarshal(jsonData, &cr.FilteredMessages)
	if err != nil {
		log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: json.Unmarshal failed")
		return err
	}
	return nil
}

func AppendFieldToSchema(f artemis_orchestrations.JsonSchemaField, jd artemis_orchestrations.JsonSchemaDefinition) artemis_orchestrations.JsonSchemaDefinition {
	for i, v := range jd.Fields {
		if v.FieldName == f.FieldName {
			jd.Fields[i] = f
			return jd
		}
	}
	jd.Fields = append(jd.Fields, f)
	return jd
}
