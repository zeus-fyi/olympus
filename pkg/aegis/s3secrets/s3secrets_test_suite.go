package s3secrets

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type S3SecretsTestSuite struct {
	test_suites.S3TestSuite
	e test_suites.EncryptionTestSuite

	S3Secrets S3Secrets
}

func (t *S3SecretsTestSuite) SetupTest() {
	t.SetupLocalDigitalOceanS3()
	t.e.SetupLocalAge()

	c := compression.NewCompression()
	r := s3reader.NewS3ClientReader(t.S3)

	t.S3Secrets = NewS3Secrets(c, t.e.Age, r)
}

func (t *S3SecretsTestSuite) TestReadGzipAndEncryptDecrypt() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("test.txt"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		Fn:          "unencrypted-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	reader := s3reader.NewS3ClientReader(t.S3)
	err := reader.Read(ctx, p, input)
	t.Require().Nil(err)

	err = t.e.Age.Encrypt(&p)
	t.Require().Nil(err)

	p.Fn = "unencrypted-text.txt.age"
	p.FnOut = "decrypted-text.txt"
	err = t.e.Age.Decrypt(&p)
	t.Require().Nil(err)
}

func TestS3SecretsTestSuite(t *testing.T) {
	suite.Run(t, new(S3SecretsTestSuite))
}
