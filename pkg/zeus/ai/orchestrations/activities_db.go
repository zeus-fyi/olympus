package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/cvcio/twitter"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/openai"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
)

func (z *ZeusAiPlatformActivities) InsertEmailIfNew(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("SaveNewEmail: failed")
		return 0, err
	}
	return emailID, nil
}

func (z *ZeusAiPlatformActivities) UpsertAiOrchestration(ctx context.Context, ou org_users.OrgUser, wfParentID string, wfExecParams artemis_orchestrations.WorkflowExecParams) (int, error) {
	id, err := artemis_orchestrations.UpsertAiOrchestration(ctx, ou, wfParentID, wfExecParams)
	if err != nil {
		log.Err(err).Msg("UpsertAiOrchestration: activity failed")
		return id, err
	}
	return id, nil
}

func (z *ZeusAiPlatformActivities) InsertTelegramMessageIfNew(ctx context.Context, ou org_users.OrgUser, msg hera_search.TelegramMessage) (int, error) {
	tgId, err := hera_search.InsertNewTgMessages(ctx, ou, msg)
	if err != nil {
		log.Err(err).Interface("msg", msg).Msg("InsertTelegramMessageIfNew: failed")
		return 0, err
	}
	return tgId, nil
}

func (z *ZeusAiPlatformActivities) InsertAiResponse(ctx context.Context, msg hermes_email_notifications.EmailContents) (int, error) {
	emailID, err := hera_openai_dbmodels.InsertNewEmails(ctx, msg)
	if err != nil {
		log.Err(err).Msg("InsertAiResponse: InsertNewEmails: failed")
		return 0, err
	}
	return emailID, nil
}

func (z *ZeusAiPlatformActivities) SelectTwitterSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) (*hera_search.TwitterSearchQuery, error) {
	sq, err := hera_search.SelectTwitterSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectTwitterSearchQuery")
		return nil, err
	}
	if sq == nil {
		return nil, fmt.Errorf("SelectTwitterSearchQuery: sq is nil")
	}
	return sq, nil
}

func (z *ZeusAiPlatformActivities) SelectDiscordSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) (*hera_search.DiscordSearchResultWrapper, error) {
	sq, err := hera_search.SelectDiscordSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectDiscordSearchQuery")
		return nil, err
	}
	return sq, nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingTweetsFromSearch(ctx context.Context, searchID int, tweets []*twitter.Tweet) error {
	_, err := hera_search.InsertIncomingTweets(ctx, searchID, tweets)
	if err != nil {
		log.Err(err).Msg("InsertIncomingTweetsFromSearch")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingRedditDataFromSearch(ctx context.Context, searchID int, redditData []*reddit.Post) error {
	_, err := hera_search.InsertIncomingRedditPosts(ctx, searchID, redditData)
	if err != nil {
		log.Err(err).Msg("InsertIncomingRedditDataFromSearch")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SelectRedditSearchQuery(ctx context.Context, ou org_users.OrgUser, groupName string) ([]*hera_search.RedditSearchQuery, error) {
	rs, err := hera_search.SelectRedditSearchQuery(ctx, ou, groupName)
	if err != nil {
		log.Err(err).Msg("SelectRedditSearchQuery: activity failed")
		return nil, err
	}
	return rs, nil
}

func (z *ZeusAiPlatformActivities) InsertIncomingDiscordDataFromSearch(ctx context.Context, searchID int, messages hera_discord.ChannelMessages) error {
	err := hera_search.InsertDiscordGuild(ctx, messages.Guild.Id, messages.Guild.Name)
	if err != nil {
		log.Err(err).Msg("InsertIncomingDiscordDataFromSearch")
		return err
	}
	err = hera_search.InsertDiscordChannel(ctx, searchID, messages.Guild.Id, messages.Channel.Id, messages.Channel.CategoryId, messages.Channel.Category, messages.Channel.Name, messages.Channel.Topic)
	if err != nil {
		log.Err(err).Msg("InsertIncomingDiscordDataFromSearch")
		return err
	}
	_, err = hera_search.InsertIncomingDiscordMessages(ctx, searchID, messages)
	if err != nil {
		log.Err(err).Interface("msgs", messages.Messages).Msg("InsertIncomingRedditDataFromSearch")
		return err
	}
	return nil
}
