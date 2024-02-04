package hestia_login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_reddit "github.com/zeus-fyi/olympus/pkg/hera/reddit"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func RedditLoginHandler(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, Conf.AuthCodeURL(state))
}

func RedditCallbackHandler(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(c.Request().Context(), ou, "reddit")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	rc, err := hera_reddit.InitOrgRedditClient(c.Request().Context(), ps.OAuth2Public, ps.OAuth2Secret, ps.Username, ps.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	meInfo, err := rc.GetMe(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if meInfo == nil || meInfo.Name == "" {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	sr := SecretsRequest{
		Name:  fmt.Sprintf("api-reddit-%s", meInfo.Name),
		Key:   "mockingbird",
		Value: rc.AccessToken,
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		OrgGroupName: fmt.Sprintf("reddit-%s", meInfo.Name),
		Routes:       []string{"https://oauth.reddit.com"},
	}
	ctx := context.Background()
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: CreateOrgGroupRoutesRequest")
		return err
	}
	c.Set("orgUser", ou)
	return sr.CreateOrUpdateSecret(c, false)
}
