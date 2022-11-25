package poseidon_chain_snapshots

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	poseidon_pkg "github.com/zeus-fyi/olympus/poseidon/pkg"
)

type Response struct {
	URL string
}

func (t *DownloadChainSnapshotRequest) GeneratePresignedURL(c echo.Context) error {
	s3pos := poseidon_pkg.PoseidonReader
	ctx := context.Background()
	input := BucketFinder(t.BucketRequest)
	url, err := s3pos.GeneratePresignedURL(ctx, input)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GeneratePresignedURL")
		return err
	}
	resp := Response{URL: url}
	return c.JSON(http.StatusOK, resp)
}

func BucketFinder(br poseidon.BucketRequest) *s3.GetObjectInput {
	input := &s3.GetObjectInput{}
	switch br.ClientName {
	case "lighthouse":
		br = poseidon_buckets.LighthouseMainnetBucket
	case "geth":
		br = poseidon_buckets.GethMainnetBucket
	default:
		return input
	}
	input.Bucket = aws.String(br.BucketName)
	input.Key = aws.String(br.GetBucketKey())
	return input
}
