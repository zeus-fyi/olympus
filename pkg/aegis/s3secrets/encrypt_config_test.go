package s3secrets

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type S3SecretsEncryptTestSuite struct {
	S3SecretsTestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3SecretsEncryptTestSuite) TestGzipAndEncrypt() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		Fn:          "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := t.S3Secrets.GzipAndEncrypt(&p)
	t.Require().Nil(err)
}

func TestS3SecretsEncryptTestSuite(t *testing.T) {
	suite.Run(t, new(S3SecretsEncryptTestSuite))
}
