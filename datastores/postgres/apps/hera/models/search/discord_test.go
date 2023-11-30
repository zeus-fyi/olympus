package hera_search

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (s *SearchAITestSuite) TestSelectDiscordSearchMessagesQuery() {
	// Initialize context and necessary data
	// Setup context and necessary data
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	si := TimeInterval{}
	si[0] = time.Now().AddDate(0, 0, -7)

	fmt.Println(si[0].Unix())
	si[1] = time.Now()
	fmt.Println(si[1].Unix())

	// Call the function
	sp := AiSearchParams{
		SearchContentText:    "node",
		GroupFilter:          "",
		Platforms:            "",
		Usernames:            "",
		WorkflowInstructions: "",
		SearchInterval:       si,
		AnalysisInterval:     TimeInterval{},
	}
	results, err := SearchDiscord(ctx, ou, sp)

	// Assert expected outcomes
	s.Require().NoError(err, "SelectDiscordSearchQuery should not return an error")
	s.Require().NotNil(results, "Results should not be nil")
}

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

	f := filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/datastores/postgres/apps/hera/models/search",
		DirOut:      "",
		FnIn:        "eth-discord.json",
		FnOut:       "",
		Env:         "",
		FilterFiles: nil,
	}
	b := f.ReadFileInPath()
	messages := hera_discord.ChannelMessages{}
	err := json.Unmarshal(b, &messages)
	s.Require().NoError(err)

	err = InsertDiscordGuild(ctx, messages.Guild.Id, messages.Guild.Name)

	s.Require().NoError(err)
	err = InsertDiscordChannel(ctx, 1700781280741432832, messages.Guild.Id, messages.Channel.Id, messages.Channel.CategoryId, messages.Channel.Category, messages.Channel.Name, messages.Channel.Topic)
	s.Require().NoError(err)

	//messages := hera_discord.ChannelMessages{
	//	Guild: hera_discord.Guild{
	//		Id:   "exampleGuildID",
	//		Name: "exampleGuildID",
	//	},
	//	Channel: hera_discord.Channel{
	//		Id:         "testChannelID",
	//		CategoryId: "123",
	//		Category:   "Example Category",
	//		Name:       "Example Channel",
	//		Topic:      "Sample Topic",
	//	},
	//	Messages: []hera_discord.Message{
	//		{
	//			Author: hera_discord.Author{
	//				Id:       "author_1",
	//				Name:     "Author One",
	//				Nickname: "Author1Nick",
	//				Roles: []hera_discord.Role{
	//					{Id: "role1", Name: "Role One"},
	//				},
	//			},
	//			Content:  "This is a test message",
	//			Id:       "1700781280741432832",
	//			Mentions: []hera_discord.Mention{},
	//			Reactions: []hera_discord.Reaction{
	//				{
	//					Count: 1,
	//					Emoji: hera_discord.Emoji{
	//						Code: "emoji1",
	//					},
	//				},
	//			},
	//			TimestampEdited: time.Now(), // use a specific time if necessary
	//			Type:            "messageType",
	//			Reference: hera_discord.Reference{
	//				ChannelId: "ref_channel_1",
	//				GuildId:   "ref_guild_1",
	//				MessageId: "ref_msg_1",
	//			},
	//		},
	//	},
	//	MessageCount: 1, // Set this to the actual number of messages in Messages
	//}

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
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
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