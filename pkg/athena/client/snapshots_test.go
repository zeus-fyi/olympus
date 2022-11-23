package athena_client

import (
	"github.com/zeus-fyi/olympus/pkg/athena/client/poseidon_buckets"
)

func (t *AthenaClientTestSuite) UploadTest() {
	//br := poseidon_buckets.GethMainnetBucket
	br := poseidon_buckets.LighthouseMainnetBucket
	err := t.AthenaTestClient.Upload(ctx, br)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) DownloadTest() {
	//br := poseidon_buckets.GethMainnetBucket
	br := poseidon_buckets.LighthouseMainnetBucket
	err := t.AthenaTestClient.Download(ctx, br)
	t.Assert().Nil(err)
}
