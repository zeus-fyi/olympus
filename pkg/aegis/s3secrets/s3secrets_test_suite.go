package s3secrets

import (
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type S3SecretsManagerTestSuite struct {
	test_suites_s3.S3TestSuite
	e test_suites_encryption.EncryptionTestSuite

	S3Secrets S3Secrets
}

func (t *S3SecretsManagerTestSuite) SetupTest() {
	t.SetupLocalDigitalOceanS3()
	t.e.SetupLocalAge()

	r := s3reader.NewS3ClientReader(t.S3)
	t.S3Secrets = NewS3Secrets(t.e.Age, r)
}
