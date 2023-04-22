package hera_openai_dbmodels

import (
	"context"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type HeraOpenAITestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *HeraOpenAITestSuite) TestSelectBalance() {
	ctx := context.Background()
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	b, err := CheckTokenBalance(ctx, ou)

	s.Require().Nil(err)
	s.Assert().NotZero(b.TokensConsumed)
	s.Assert().NotZero(b.TokensRemaining)
}

func (s *HeraOpenAITestSuite) TestInsertCompletionResponse() {
	s.InitLocalConfigs()
	ctx := context.Background()
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	response := openai.CompletionResponse{
		ID:      "cmpl-GERzeJQ4lvqPk8SkZu4XMIuR",
		Object:  "text_completion",
		Created: 1586839808,
		Model:   "text-davinci:003",
		Choices: []openai.CompletionChoice{
			{
				Text:         "\n\nThis is indeed a test",
				Index:        0,
				FinishReason: "length",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     5,
			CompletionTokens: 7,
			TotalTokens:      12,
		},
	}
	err := InsertCompletionResponse(ctx, ou, response)
	s.Require().Nil(err)
}

func TestHeraOpenAITestSuite(t *testing.T) {
	suite.Run(t, new(HeraOpenAITestSuite))
}
