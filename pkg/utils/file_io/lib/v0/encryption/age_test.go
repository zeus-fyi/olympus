package encryption

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type AgeEncryptionTestSuite struct {
	test_suites.S3TestSuite
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

	pubKey := s.Tc.LocalAgePubkey
	err := Encrypt(p, pubKey)
	s.Require().Nil(err)
}

// use age-keygen -o private_key.txt to create a pubkey/private key pair for here
func (s *AgeEncryptionTestSuite) TestDecryption() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "",
		Fn:          "kube.tar.gz.age",
		FnOut:       "kube_decrypted.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := Decrypt(p, s.Tc.LocalAgePkey)
	s.Require().Nil(err)
}

func TestAgeEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(AgeEncryptionTestSuite))
}
