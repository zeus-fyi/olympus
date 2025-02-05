package athena_client

import (
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
)

func (t *AthenaClientTestSuite) TestUploadConsensusClient() {
	br := poseidon_buckets.LighthouseMainnetBucket
	err := t.AthenaTestClient.Upload(ctx, br)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) TestUploadExecClient() {
	br := poseidon_buckets.GethMainnetBucket
	err := t.AthenaTestClient.Upload(ctx, br)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) TestDownloadConsensusClient() {
	br := poseidon_buckets.LighthouseMainnetBucket
	err := t.AthenaTestClient.Download(ctx, br)
	t.Assert().Nil(err)
}
