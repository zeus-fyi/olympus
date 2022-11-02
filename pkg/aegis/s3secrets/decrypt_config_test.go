package s3secrets

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type S3SecretsDecryptTestSuite struct {
	S3SecretsTestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3SecretsDecryptTestSuite) TestReadGzipAndEncryptDecrypt() {

}

func TestS3SecretsDecryptTestSuite(t *testing.T) {
	suite.Run(t, new(S3SecretsDecryptTestSuite))
}
