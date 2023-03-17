package artemis_api_router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/artemis/api/v1/ethereum/send_tx"
	artemis_ethereum_validator_service "github.com/zeus-fyi/olympus/artemis/api/v1/ethereum/validator_service"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)

	//go func() {
	//	artemis_eth_txs.Faucet = artemis_eth_txs.NewFaucetServer()
	//	artemis_eth_txs.Faucet.Run()
	//}()

	//e.POST("/v1beta/ethereum/ephemery/send/api/claim", artemis_eth_txs.SendEtherEphemeralFaucetTxHandler)
	//e.POST("/ethereum/ephemeral/send/api/info",s.handleInfo()))

	InitV1BetaRoutes(e)
	InitV1Routes(e)
	InitV1InternalRoutes(e)
	return e
}

func InitV1BetaRoutes(e *echo.Echo) {
	eg := e.Group("/v1beta")
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

	eg.POST("/ethereum/goerli/tx", artemis_eth_txs.SendSignedTxEthGoerliTxHandler)
	eg.POST("/ethereum/goerli/send", artemis_eth_txs.SendEtherGoerliTxHandler)
	eg.POST("/ethereum/ephemery/send", artemis_eth_txs.SendEtherEphemeralTxHandler)
	eg.POST("/ethereum/ephemery/tx", artemis_eth_txs.SendSignedTxEthEphemeralTxHandler)
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
	eg.POST("/ethereum/validators/create", artemis_ethereum_validator_service.CreateEthereumValidatorsHandler)
}

func InitV1InternalRoutes(e *echo.Echo) {
	eg := e.Group("/v1/internal")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyInternalBearerToken(ctx, token)
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
}
func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
