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
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	age_encryption "github.com/zeus-fyi/zeus/pkg/aegis/crypto/age"
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
			cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
			if err == nil && cookie != nil {
				log.Info().Msg("InitV1Routes: Cookie found")
				token = cookie.Value
			}
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

	eg.GET("/age/generate", GenerateRandomAgeEncryptionKey) // if no js client, generate age keypair

	// ethereum aws automation
	// lambda user & role policy creation
	eg.POST("/ethereum/validators/aws/user/internal/lambda/create", v1_ethereum_aws.CreateServerlessInternalUserHandler)
	eg.POST("/ethereum/validators/aws/user/external/lambda/create", v1_ethereum_aws.CreateServerlessExternalUserHandler)
	// lambda access keys creation
	eg.POST("/ethereum/validators/aws/lambda/external/user/access/keys/create", v1_ethereum_aws.CreateServerlessExternalUserAuthHandler)

	// lambda signers
	eg.POST("/ethereum/validators/aws/lambda/signer/create", v1_ethereum_aws.CreateBlsLambdaFunctionHandler)
	eg.POST("/ethereum/validators/aws/lambda/signer/keystores/layer/create", v1_ethereum_aws.CreateServerlessKeystoresLayerHandler)

	// lambda create
	eg.POST("/ethereum/validators/aws/lambda/secrets/create", v1_ethereum_aws.CreateLambdaFunctionSecretsKeyGenHandler)
	eg.POST("/ethereum/validators/aws/lambda/keystores/zip/create", v1_ethereum_aws.CreateLambdaFunctionEncZipGenHandler)
	eg.POST("/ethereum/validators/aws/lambda/deposits/create", v1_ethereum_aws.CreateLambdaFunctionDepositsGenHandler)

	// lambda verify
	eg.POST("/ethereum/validators/aws/lambda/url", v1_ethereum_aws.GetLambdaFunctionURLHandler)
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

type GeneratedAgeKey struct {
	AgePrivateKey string `json:"agePrivateKey"`
	AgePublicKey  string `json:"agePublicKey"`
}

func GenerateRandomAgeEncryptionKey(c echo.Context) error {
	key := GeneratedAgeKey{}
	key.AgePublicKey, key.AgePrivateKey = age_encryption.GenerateNewKeyPair()
	return c.JSON(http.StatusOK, key)
}
