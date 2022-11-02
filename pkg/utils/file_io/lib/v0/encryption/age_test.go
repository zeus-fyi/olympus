package encryption

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type AgeEncryptionTestSuite struct {
	base.CoreTestSuite
}

func (s *AgeEncryptionTestSuite) TestEncryption() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		Fn:          "kube.tar.gz",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	pubKey := "age1f2awqn4xvrp4sehrv6zq0s64lt278hh7vq6darny4kzmlhfnusxq3hf62a"
	err := Encrypt(p, pubKey)
	s.Require().Nil(err)
}

func (s *AgeEncryptionTestSuite) TestDecryption() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		Fn:          "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := Decrypt(p, "")
	s.Require().Nil(err)
}

func TestAgeEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(AgeEncryptionTestSuite))
}
