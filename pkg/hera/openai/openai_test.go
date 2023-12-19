package hera_openai

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-resty/resty/v2"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type HeraTestSuite struct {
	test_suites_base.TestSuite
}

func (s *HeraTestSuite) TestOpenAIGetModels() {
	r := resty.New()
	r.SetAuthToken(s.Tc.OpenAIAuth)
	resp, err := r.R().
		Get("https://api.openai.com/v1/models")
	s.Require().Nil(err)
	fmt.Println(string(resp.Body()))

}

func (s *HeraTestSuite) TestOpenAICreateJsonOutputFormat() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	fdSchema := jsonschema.Definition{
		Type:        jsonschema.Object,
		Description: "",
		Enum:        nil,
		Properties: map[string]jsonschema.Definition{
			"location": {
				Type:        jsonschema.String,
				Description: "The city and state, e.g. San Francisco, CA",
			},
			"unit": {
				Type: jsonschema.String,
				Enum: []string{"celcius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
		Items:    nil,
	}
	fd := openai.FunctionDefinition{
		Name:        "fn-json",
		Description: "",
		Parameters:  fdSchema,
	}
	params := OpenAIParams{
		Model:              "gpt-4-1106-preview",
		MaxTokens:          300,
		Prompt:             "what is the meaning of life",
		FunctionDefinition: fd,
	}
	resp, err := HeraOpenAI.MakeCodeGenRequestJsonFormattedOutput(context.Background(), ou, params)
	s.Require().Nil(err)
	s.Require().NotEmpty(resp)
}

func (s *HeraTestSuite) TestOpenAICreateAssistant() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	InitHeraOpenAI(s.Tc.OpenAIAuth)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	resp, err := HeraOpenAI.CreateAssistant(
		context.Background(),
		openai.AssistantRequest{
			Model:        "gpt-4-1106-preview",
			Name:         aws.String(fmt.Sprintf("%d", ou.UserID)),
			Description:  aws.String(fmt.Sprintf("%s", "description")),
			Instructions: aws.String(fmt.Sprintf("%s", "instructions")),
			Tools: []openai.AssistantTool{
				{
					Type: "",
					Function: &openai.FunctionDefinition{
						Name:        "",
						Description: "",
						Parameters:  nil,
					},
				},
			},
			FileIDs:  nil,
			Metadata: nil,
		},
	)
	s.Require().Nil(err)
	fmt.Println(resp)
}

func (s *HeraTestSuite) TestOpenAIChatGptInsert() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	InitHeraOpenAI(s.Tc.OpenAIAuth)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	resp, err := HeraOpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "write me a go for loop and be brief",
					Name:    fmt.Sprintf("%d", ou.UserID),
				},
			},
		},
	)
	s.Require().Nil(err)
	fmt.Println(resp)
	err = HeraOpenAI.RecordUIChatRequestUsage(ctx, ou, resp, nil)
	s.Require().Nil(err)
}

func (s *HeraTestSuite) TestOpenAITokenCount() {
	ForceDirToPythonDir()
	bytes, err := os.ReadFile("./example.txt")
	s.Require().Nil(err)
	tokenCount := GetTokenApproximate(string(bytes))
	s.Assert().Equal(61, tokenCount)
	// NOTE open gpt-3 https://beta.openai.com/tokenizer returns 64 tokens as the count
	// there's no opensource transformer for this, so use this + some margin when sending requests
	// 2048 is the max token count for most models, the max size - prompt size, is your limitation on completion
	// tokens
}

func (s *HeraTestSuite) TestOpenAI() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	InitHeraOpenAI(s.Tc.OpenAIAuth)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	model := gogpt.GPT3TextDavinci003

	params := OpenAIParams{
		Model:     model,
		MaxTokens: 300,
		Prompt:    "what is the meaning of life",
	}
	resp, err := HeraOpenAI.MakeCodeGenRequest(ctx, ou, params)
	s.Require().Nil(err)
	fmt.Println(resp)

}

func TestHeraTestSuite(t *testing.T) {
	suite.Run(t, new(HeraTestSuite))
}
