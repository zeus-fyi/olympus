package hera_openai_dbmodels

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

type TestTelegramMetadata struct {
	IsReply       bool   `json:"is_reply,omitempty"`
	IsChannel     bool   `json:"is_channel,omitempty"`
	IsGroup       bool   `json:"is_group,omitempty"`
	IsPrivate     bool   `json:"is_private,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	MutualContact bool   `json:"mutual_contact,omitempty"`
	Username      string `json:"username,omitempty"`
}

func (s *HeraOpenAITestSuite) TestInsertTelegramMsg() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
	ctx := context.Background()
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	tm := TestTelegramMetadata{
		IsReply:       false,
		IsChannel:     false,
		IsGroup:       false,
		IsPrivate:     false,
		FirstName:     "fir",
		LastName:      "sdf",
		Phone:         "6267282",
		MutualContact: false,
		Username:      "sadfuoasd",
	}
	b, err := json.Marshal(tm)
	s.Require().Nil(err)

	re, err := InsertNewTgMessages(ctx, ou, 1586839808, 123456789, 123456789, 123456789, "test", "testsdfsd", b)
	s.Require().Nil(err)
	s.Assert().NotZero(re)
}

func TestHeraOpenAITestSuite(t *testing.T) {
	suite.Run(t, new(HeraOpenAITestSuite))
}
