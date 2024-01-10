package artemis_orchestrations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type AiAssistant struct {
	openai.Assistant
}

// TODO, if user has a key, create it on their account, else don't create

func InsertAssistant(ctx context.Context, ou org_users.OrgUser, assistant openai.Assistant) error {
	return nil
}

// TODO, lookup users assistants via API key, then update the database if needed

func SelectAssistants(ctx context.Context, ou org_users.OrgUser) ([]AiAssistant, error) {
	tmp := []AiAssistant{
		{
			Assistant: openai.Assistant{
				ID:           "test-id1",
				Name:         aws.String("test-name"),
				Model:        "gpt-3.5",
				Instructions: aws.String("This is a test assistant."),
			},
		},
		{
			Assistant: openai.Assistant{
				ID:           "test-id2",
				Model:        "gpt-4",
				Instructions: aws.String("This is a test assistant."),
			},
		},
	}
	return tmp, nil
}
