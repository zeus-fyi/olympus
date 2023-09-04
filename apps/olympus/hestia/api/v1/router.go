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
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/hestia/api/v1/ethereum/aws"
	hestia_iris_v1_routes "github.com/zeus-fyi/olympus/hestia/api/v1/iris"
	hestia_quiknode_v1_routes "github.com/zeus-fyi/olympus/hestia/api/v1/quiknode"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	age_encryption "github.com/zeus-fyi/zeus/pkg/aegis/crypto/age"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/healthz", Health)
	e.GET("/healthcheck", Health)
	e.GET("/health", Health)

	hestia_quiknode_v1_routes.InitV1RoutesServices(e)
	InitV1Routes(e)
	InitV1InternalRoutes(e)
	return e
}

const (
	IrisCreateRoutesPath           = "/iris/routes/create"
	IrisReadRoutesPath             = "/iris/routes/read"
	IrisReadAllRoutesAndGroupsPath = "/iris/routes/read/all"

	IrisDeleteRoutesPath = "/iris/routes/delete"

	IrisReadGroupRoutesPath                = "/iris/routes/group/:groupName/read"
	IrisReadGroupTableMetricsPath          = "/iris/routes/group/:groupName/metrics"
	IrisUpdateGroupRoutesPath              = "/iris/routes/group/:groupName/update"
	IrisRemoveEndpointsFromGroupRoutesPath = "/iris/routes/group/:groupName/endpoints/delete"

	IrisCreateGroupRoutesPath = "/iris/routes/groups/create"
	IrisReadGroupsRoutesPath  = "/iris/routes/groups/read"
	IrisDeleteGroupRoutesPath = "/iris/routes/groups/delete"

	IrisDeleteRoutesPathInternal = "/iris/routes/delete/:orgID"
)

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
			k := read_keys.NewKeyReader()
			services, err := k.QueryUserAuthedServices(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
				return false, err
			}
			if len(services) <= 0 {
				log.Warn().Msg("InitV1Routes: No services found")
				return false, nil
			}
			c.Set("servicePlans", k.Services)
			ou := org_users.NewOrgUserWithID(k.OrgID, k.GetUserID())
			c.Set("orgUser", ou)
			c.Set("bearer", k.PublicKey)
			return len(services) > 0, err
		},
	}))

	eg.POST(IrisCreateRoutesPath, hestia_iris_v1_routes.CreateOrgRoutesRequestHandler)
	eg.GET(IrisReadRoutesPath, hestia_iris_v1_routes.ReadOrgRoutesRequestHandler)
	eg.DELETE(IrisDeleteRoutesPath, hestia_iris_v1_routes.DeleteOrgRoutesRequestHandler)

	eg.GET(IrisReadGroupRoutesPath, hestia_iris_v1_routes.ReadOrgGroupRoutesRequestHandler)
	eg.GET(IrisReadGroupsRoutesPath, hestia_iris_v1_routes.ReadOrgGroupsRoutesRequestHandler)
	eg.GET(IrisReadGroupTableMetricsPath, hestia_iris_v1_routes.ReadTableMetricsRequestHandler)

	eg.POST(IrisCreateGroupRoutesPath, hestia_iris_v1_routes.CreateOrgGroupRoutesRequestHandler)
	eg.PUT(IrisUpdateGroupRoutesPath, hestia_iris_v1_routes.UpdateOrgGroupRoutesRequestHandler)
	eg.DELETE(IrisDeleteGroupRoutesPath, hestia_iris_v1_routes.DeletePartialOrgGroupRoutesRequestHandler)
	eg.GET(IrisReadAllRoutesAndGroupsPath, hestia_iris_v1_routes.ReadAllOrgGroupsAndEndpointsRequestHandler)

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
	eg.DELETE(IrisDeleteRoutesPathInternal, hestia_iris_v1_routes.InternalDeleteOrgRoutesRequestHandler)

	//eg.POST("/users/create", CreateUserHandler)
	//eg.POST("/orgs/create", CreateOrgHandler)
	//eg.POST("/cloud/namespace/request/create", CreateTopologiesOrgCloudCtxNsRequestHandler)
}

func Health(c echo.Context) error {
	resp := Response{Message: "ok"}
	return c.JSON(http.StatusOK, resp)
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
