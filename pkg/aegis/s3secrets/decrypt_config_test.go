package s3secrets

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type S3SecretsDecryptTestSuite struct {
	S3SecretsManagerTestSuite
}

func (t *S3SecretsDecryptTestSuite) TestPullAndGzipAndDecryptToInMemFs() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("kube.tar.gz.age"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "kube.tar.gz.age",
		FnOut:       "kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	unzipDir := "./.kube"
	err := t.S3Secrets.PullS3AndDecryptAndUnGzipToInMemFs(ctx, &p, unzipDir, input)
	t.Require().Nil(err)
}

func (t *S3SecretsDecryptTestSuite) TestDecryptAndUnGzipInMemFs() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "kube.tar.gz.age",
		FnOut:       "kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	unzipDir := "./kube"
	err := t.S3Secrets.DecryptAndUnGzipToInMemFs(&p, unzipDir)
	t.Require().Nil(err)
}

// TestRead, you'll need to set the secret values to run the test
func (t *S3SecretsDecryptTestSuite) TestDecryptAndUnGzip() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		FnIn:        "kube.tar.gz.age",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := t.S3Secrets.DecryptAndUnGzip(&p)
	t.Require().Nil(err)
}

func (t *S3SecretsDecryptTestSuite) TestReadGzipAndEncryptDecrypt() {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("test.txt"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		FnIn:        "unencrypted-text.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	reader := s3reader.NewS3ClientReader(t.S3)
	err := reader.Read(ctx, &p, input)
	t.Require().Nil(err)

	err = t.e.Age.Encrypt(&p)
	t.Require().Nil(err)

	p.FnIn = "unencrypted-text.txt.age"
	p.FnOut = "decrypted-text.txt"
	err = t.e.Age.DecryptToFile(&p)
	t.Require().Nil(err)
}

func TestS3SecretsDecryptTestSuite(t *testing.T) {
	suite.Run(t, new(S3SecretsDecryptTestSuite))
}
