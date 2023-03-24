package v1_poseidon

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)
	InitV1InternalRoutes(e)
	InitV1Routes(e)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
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
	eg.GET("/ethereum/mainnet/geth", GetGethPresignedURL)
	eg.GET("/ethereum/mainnet/lighthouse", GetLighthousePresignedURL)

	eg.POST("/ethereum/beacon/disk/wipe", DiskWipeRequestHandler)
	eg.POST("/ethereum/beacon/disk/upload", SnapshotUploadRequestHandler)
	eg.POST("/ethereum/beacon/chainsync/upload", BeaconChainSyncUploadRequestHandler)

}
func GetGethPresignedURL(c echo.Context) error {
	ctx := context.Background()
	reader := s3reader.NewS3ClientReader(poseidon_orchestrations.PoseidonS3Manager)

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi-snapshots"),
		Key:    aws.String(poseidon_buckets.GethMainnetBucket.GetBucketKey()),
	}

	url, err := reader.GeneratePresignedURL(ctx, input)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetGethPresignedURL")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, url)
}

func GetLighthousePresignedURL(c echo.Context) error {
	ctx := context.Background()
	reader := s3reader.NewS3ClientReader(poseidon_orchestrations.PoseidonS3Manager)
	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi-snapshots"),
		Key:    aws.String(poseidon_buckets.LighthouseMainnetBucket.GetBucketKey()),
	}
	url, err := reader.GeneratePresignedURL(ctx, input)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetLighthousePresignedURL")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, url)
}

func InitV1InternalRoutes(e *echo.Echo) {
	eg := e.Group("/v1beta/internal")
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
	eg.GET("/ethereum", RunWorkflow)
}

func RunWorkflow(c echo.Context) error {
	return ExecuteSyncWorkflow(c, context.Background())
}

func ExecuteSyncWorkflow(c echo.Context, ctx context.Context) error {
	err := poseidon_orchestrations.PoseidonSyncWorker.ExecutePoseidonSyncWorkflow(ctx)
	if err != nil {
		log.Err(err).Msg("ExecuteSyncWorkflow")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
