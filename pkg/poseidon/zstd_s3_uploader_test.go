package poseidon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type S3UploaderTestSuite struct {
	test_suites_s3.S3TestSuite
}

var brUploadS3 = S3BucketRequest{
	BucketName: "flows",
}

func (s *S3UploaderTestSuite) TestOvHTextFileZstdCmpAndUpload() {
	ctx := context.Background()
	pos := NewS3Poseidon(s.OvhS3)
	pos.DirIn = "/Users/alex/go/Olympus/olympus/pkg/poseidon/"
	pos.FnIn = "tmp.txt"
	err := pos.S3ZstdCompressAndUpload(ctx, brUploadS3)
	s.Require().Nil(err)
}

func TestS3UploaderTestSuite(t *testing.T) {
	suite.Run(t, new(S3UploaderTestSuite))
}
