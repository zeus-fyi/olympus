package hera_search

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestInsertDiscordSearchQuery() {
	// Setup context and necessary data
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	searchGroupName := "exampleGroupName"
	maxResults := 100
	query := "exampleQuery"

	// Call the function
	searchID, err := InsertDiscordSearchQuery(ctx, ou, searchGroupName, maxResults, query)

	// Assert expected outcomes
	s.Require().NoError(err)
	s.Assert().NotZero(searchID)
	fmt.Println(searchID)
}

func (s *SearchAITestSuite) TestInsertDiscordChannel() {
	// Initialize context and necessary data
	searchID := 1700781280741432832 // Replace with a valid search ID
	guildID := "exampleGuildID"
	channelID := "testChannelID"
	categoryID := "testCategoryID"
	category := "testCategory"
	name := "testChannelName"
	topic := "testTopic"

	// Call the function
	err := InsertDiscordChannel(ctx, searchID, guildID, channelID, categoryID, category, name, topic)

	// Assert expected outcomes
	s.Require().NoError(err, "InsertDiscordChannel should not return an error")

	// Further assertions can be made here to verify the state of the database
	// This might include querying the database to ensure the record was inserted/updated correctly
}

func (s *SearchAITestSuite) TestInsertDiscordGuild() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// Setup context and necessary data
	guildID := "exampleGuildID"
	name := "exampleGuildName"

	// Call the function
	err := InsertDiscordGuild(ctx, guildID, name)

	// Assert expected outcomes
	s.Require().NoError(err)
}

func (s *SearchAITestSuite) TestInsertIncomingDiscordMessages() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	// Initialize context and necessary data
	messages := []*DiscordMessage{
		{
			MessageID: 1700781280741432831,
			SearchId:  1700781280741432832,
			GuildID:   "exampleGuildID",
			ChannelID: "testChannelID",
			Author:    json.RawMessage(`{"name":"Author1"}`),
			Content:   "Message content 1",
			Mentions:  json.RawMessage(`[{"id":"user1"}]`),
			Reactions: json.RawMessage(`[{"count":5,"emoji":{"code":"emoji1"}}]`),
			Reference: json.RawMessage(`{}`),
			EditedAt:  0,
			Type:      "messageType1",
		},
		// Add more mock messages if needed
	}

	// Call the function
	messageIDs, err := InsertIncomingDiscordMessages(ctx, messages)

	// Assert expected outcomes
	s.Require().NoError(err, "InsertIncomingDiscordMessages should not return an error")
	s.Assert().NotNil(messageIDs, "Returned message IDs should not be nil")
	s.Assert().Len(messageIDs, len(messages), "The number of returned message IDs should match the number of input messages")

	// Additional checks can be added here, such as verifying the contents of the returned message IDs
}
func (s *SearchAITestSuite) TestSelectDiscordSearchQuery() {
	// Initialize context and necessary data
	// Setup context and necessary data
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	searchGroupName := "exampleGroupName"

	// Call the function
	results, err := SelectDiscordSearchQuery(ctx, ou, searchGroupName)

	// Assert expected outcomes
	s.Require().NoError(err, "SelectDiscordSearchQuery should not return an error")
	s.Require().NotNil(results, "Results should not be nil")

}
