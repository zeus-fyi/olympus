package zeus_v1_ai

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type AiSearchIndexerRequest struct {
	SearchIndexerParamsWrapper `json:"searchIndexer"`
	PlatformSecretReference    `json:"platformSecretReference"`
}

type SearchIndexerParamsWrapper struct {
	hera_search.SearchIndexerParams
	DiscordIndexerOpts  `json:"discordOpts,omitempty"`
	EntitiesIndexerOpts `json:"entitiesOpts,omitempty"`
}

type EntitiesIndexerOpts struct {
	Nickname       string   `json:"nickname" db:"nickname"`
	EntityPlatform string   `json:"platform" db:"platform"`
	FirstName      *string  `json:"firstName,omitempty"`
	LastName       *string  `json:"lastName,omitempty"`
	Labels         []string `json:"labels"`
}

type DiscordIndexerOpts struct {
	GuildID   string `json:"guildID"`
	ChannelID string `json:"channelID"`
}

type PlatformSecretReference struct {
	SecretGroupName string `json:"secretGroupName"`
	SecretKeyName   string `json:"secretKeyName,omitempty"`
}

func AiSearchIndexerRequestHandler(c echo.Context) error {
	request := new(AiSearchIndexerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateSearchIndex(c)
}

func (r *AiSearchIndexerRequest) CreateOrUpdateSearchIndex(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if berr != nil {
		log.Error().Err(berr).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}
	if len(r.SecretKeyName) <= 0 {
		r.SecretKeyName = r.Platform
	}
	r.MaxResults = 100
	switch r.Platform {
	case "reddit":
		resp, err := hera_search.InsertRedditSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting reddit search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "twitter":
		resp, err := hera_search.InsertTwitterSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting twitter search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "discord":
		if len(r.DiscordIndexerOpts.GuildID) <= 0 || len(r.DiscordIndexerOpts.ChannelID) <= 0 {
			return c.JSON(http.StatusBadRequest, nil)
		}
		act := ai_platform_service_orchestrations.NewZeusAiPlatformActivities()
		resp, err := hera_search.InsertDiscordSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting discord search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		tn := time.Now().Add(-time.Hour * 24 * 7).Format(time.RFC3339)
		err = act.CreateDiscordJob(c.Request().Context(), ou, resp, r.DiscordIndexerOpts.ChannelID, tn)
		if err != nil {
			log.Err(err).Msg("error creating discord job")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "telegram":
	case "entities":
		eo := r.EntitiesIndexerOpts
		evs, err := artemis_entities.SelectUserMetadataByProvidedFields(c.Request().Context(), ou, eo.Nickname, eo.EntityPlatform, eo.Labels, 0)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusOK, evs)
	}
	return c.JSON(http.StatusOK, nil)
}
