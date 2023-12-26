package ai_platform_service_orchestrations

/*

act stages

objectives:
	score the output of the function definition
	save the results of the eval

	figure out how/where to inject eval stage in analysis-aggregation chains


needed, but not yet implemented:
	saving mod

activities flow:
	1. get eval fns if any
	CreatJsonOutputModelResponse

type EvalMetric struct {
    EvalMetricID          *int    `json:"evalMetricID"`
    EvalModelPrompt       string  `json:"evalModelPrompt"`
    EvalMetricName        string  `json:"evalMetricName"`
    EvalMetricResult      string  `json:"evalMetricResult"`
    EvalComparisonBoolean *bool   `json:"evalComparisonBoolean,omitempty"`
    EvalComparisonNumber  *int    `json:"evalComparisonNumber,omitempty"`
    EvalComparisonString  *string `json:"evalComparisonString,omitempty"`
    EvalMetricDataType    string  `json:"evalMetricDataType"`
    EvalOperator          string  `json:"evalOperator"`
    EvalState             string  `json:"evalState"`
}

type OpenAIParams struct {
	Model              string                    `json:"model"`
	MaxTokens          int                       `json:"maxTokens"`
	Prompt             string                    `json:"prompt"`
	FunctionDefinition openai.FunctionDefinition `json:"functionDefinition,omitempty"`
}


	fdSchema := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"count": {
				Type:        jsonschema.Number,
				Description: "total number of words in sentence",
			},
			"words": {
				Type:        jsonschema.Array,
				Description: "list of words in sentence",
				Items: &jsonschema.Definition{
					Type: jsonschema.String,
				},
			},
		},
		Required: []string{"count", "words"},
	}

	fd := openai.FunctionDefinition{
		Name:       "test",
		Parameters: fdSchema,
	}
	params := OpenAIParams{
		Model:              "gpt-4-1106-preview",
		Prompt:             "how many words are in this sentence: what is the meaning of time dilation?",
		FunctionDefinition: fd,
	}
	resp, err := HeraOpenAI.MakeCodeGenRequestJsonFormattedOutput(context.Background(), ou, params)
	s.Require().Nil(err)
	s.Require().NotEmpty(resp)
	fmt.Println(resp)

	m := map[string]interface{}{}

	for _, msg := range resp.Choices {
		for _, tool := range msg.Message.ToolCalls {
			fmt.Println(tool.Function.Name)
			err = json.Unmarshal([]byte(tool.Function.Arguments), &m)
			s.Require().Nil(err)
			count, ok := m["count"].(int)
			s.Require().True(ok)
			s.Require().Equal(7, count)
		}

	}
}
CREATE TABLE public.eval_fns(
    eval_id BIGINT PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    eval_name text NOT NULL,
    eval_type text NOT NULL,
    eval_group_name text NOT NULL,
    eval_model text,
    eval_format text NOT NULL
);

CREATE TABLE public.eval_metrics(
    eval_metric_id BIGINT PRIMARY KEY,
    eval_id BIGINT NOT NULL REFERENCES public.eval_fns(eval_id),
    eval_model_prompt text NOT NULL,
    eval_metric_name text NOT NULL,
    eval_metric_result text NOT NULL,
    eval_comparison_boolean boolean,
    eval_comparison_number BIGINT,
    eval_comparison_string text,
    eval_metric_data_type text NOT NULL,
    eval_operator text NOT NULL,
    eval_state text NOT NULL
);
*/
