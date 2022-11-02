package s3secrets

import (
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type S3SecretsManagerTestSuite struct {
	test_suites.S3TestSuite
	e test_suites.EncryptionTestSuite

	S3Secrets S3Secrets
}

func (t *S3SecretsManagerTestSuite) SetupTest() {
	t.SetupLocalDigitalOceanS3()
	t.e.SetupLocalAge()

	c := compression.NewCompression()
	r := s3reader.NewS3ClientReader(t.S3)

	t.S3Secrets = NewS3Secrets(c, t.e.Age, r)
}
