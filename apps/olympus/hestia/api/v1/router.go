package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/hestia/api/v1/ethereum/aws"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)

	InitV1Routes(e)
	InitV1InternalRoutes(e)
	return e
}

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerToken(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))

	eg.GET("/age/generate", TODO) // if no js client, generate age keypair

	// ethereum aws automation
	// validator deposit & keystore generation
	eg.POST("/ethereum/validators/aws/generation", v1_ethereum_aws.GenerateValidatorsHandler)
	// lambda related
	eg.POST("/ethereum/validators/aws/user/internal/lambda/create", v1_ethereum_aws.CreateServerlessInternalUserHandler)
	eg.POST("/ethereum/validators/aws/user/external/lambda/create", v1_ethereum_aws.CreateServerlessExternalUserHandler)
	eg.POST("/ethereum/validators/aws/lambda/keystore/create", v1_ethereum_aws.CreateServerlessKeystoresHandler)
	eg.POST("/ethereum/validators/aws/lambda/create", v1_ethereum_aws.CreateLambdaFunctionHandler)
	eg.POST("/ethereum/validators/aws/lambda/verify", v1_ethereum_aws.VerifyLambdaFunctionHandler)

	// zeus service
	eg.GET("/ethereum/validators/service/info", GetValidatorServiceInfoHandler)
	eg.POST("/ethereum/validators/service/create", CreateValidatorServiceRequestHandler)
	eg.POST("/validators/service/create", CreateValidatorServiceRequestHandler)
}

func InitV1InternalRoutes(e *echo.Echo) {
	eg := e.Group("/v1/internal")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerTokenService(ctx, token, create_org_users.ZeusWebhooksService)
			if err != nil {
				log.Err(err).Msg("InitV1InternalRoutes")
				return false, c.JSON(http.StatusInternalServerError, nil)
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", key.PublicKey)
			return key.PublicKeyVerified, err
		},
	}))
	eg.POST(DemoUsersCreateRoute, CreateDemoUserHandler)
	//eg.POST("/users/create", CreateUserHandler)
	//eg.POST("/orgs/create", CreateOrgHandler)
	//eg.POST("/cloud/namespace/request/create", CreateTopologiesOrgCloudCtxNsRequestHandler)
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

func TODO(c echo.Context) error {
	return c.String(http.StatusOK, "TODO")
}
