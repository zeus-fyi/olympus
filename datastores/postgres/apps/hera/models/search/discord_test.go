package hera_search

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
)

func (s *SearchAITestSuite) TestInsertDiscordSearchQuery() {
	// Setup context and necessary data
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	searchGroupName := "zeusfyi"
	maxResults := 100
	query := ""

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
	messages := hera_discord.ChannelMessages{
		Guild: hera_discord.Guild{
			Id:   "exampleGuildID",
			Name: "exampleGuildID",
		},
		Channel: hera_discord.Channel{
			Id:         "testChannelID",
			CategoryId: "123",
			Category:   "Example Category",
			Name:       "Example Channel",
			Topic:      "Sample Topic",
		},
		Messages: []hera_discord.Message{
			{
				Author: hera_discord.Author{
					Id:       "author_1",
					Name:     "Author One",
					Nickname: "Author1Nick",
					Roles: []hera_discord.Role{
						{Id: "role1", Name: "Role One"},
					},
				},
				Content:         "This is a test message",
				Id:              "1700781280741432832",
				Mentions:        []hera_discord.Mention{},
				Reactions:       []hera_discord.Reaction{},
				TimestampEdited: time.Now(), // use a specific time if necessary
				Type:            "messageType",
				Reference: hera_discord.Reference{
					ChannelId: "ref_channel_1",
					GuildId:   "ref_guild_1",
					MessageId: "ref_msg_1",
				},
			},
			// ... more messages as needed
		},
		MessageCount: 1, // Set this to the actual number of messages in Messages
	}

	// Call the function
	messageIDs, err := InsertIncomingDiscordMessages(ctx, 1700781280741432832, messages)

	// Assert expected outcomes
	s.Require().NoError(err, "InsertIncomingDiscordMessages should not return an error")
	s.Assert().NotNil(messageIDs, "Returned message IDs should not be nil")

	// Additional checks can be added here, such as verifying the contents of the returned message IDs
}
func (s *SearchAITestSuite) TestSelectDiscordSearchQuery() {
	// Initialize context and necessary data
	// Setup context and necessary data
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	searchGroupName := "zeusfyi"

	// Call the function
	results, err := SelectDiscordSearchQuery(ctx, ou, searchGroupName)

	// Assert expected outcomes
	s.Require().NoError(err, "SelectDiscordSearchQuery should not return an error")
	s.Require().NotNil(results, "Results should not be nil")

}
