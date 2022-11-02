package encryption

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AgeEncryptionTestSuite struct {
	base.TestSuite
	Age Age
}

func (s *AgeEncryptionTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	pubKey := s.Tc.LocalAgePubkey
	privKey := s.Tc.LocalAgePkey
	s.Age = NewAge(privKey, pubKey)
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
	err := s.Age.Encrypt(&p)
	s.Require().Nil(err)
}

// use age-keygen -o private_key.txt to create a pubkey/private key pair for here
func (s *AgeEncryptionTestSuite) TestDecryption() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		Fn:          "kube.tar.gz.age",
		FnOut:       "kube_decrypted.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	err := s.Age.DecryptToFile(&p)
	s.Require().Nil(err)
}

func TestAgeEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(AgeEncryptionTestSuite))
}
