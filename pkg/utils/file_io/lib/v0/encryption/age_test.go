package encryption

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AgeEncryptionTestSuite struct {
	test_suites_base.TestSuite
	Age Age
}

func (s *AgeEncryptionTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	pubKey := s.Tc.LocalAgePubkey
	privKey := s.Tc.LocalAgePkey
	s.Age = NewAge(privKey, pubKey)
}

func (s *AgeEncryptionTestSuite) TestEncryption() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/pkg/utils/file_io/lib/v0/encryption",
		DirOut:      "/Users/alex/go/Olympus/olympus/pkg/utils/file_io/lib/v0/encryption",
		FnIn:        "cmds.txt",
		FnOut:       "cmds.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	err := s.Age.Encrypt(&p)
	s.Require().Nil(err)
}

func (s *AgeEncryptionTestSuite) TestItemEncryption() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./secrets",
		DirOut:      "./secrets",
		FnIn:        "key.txt",
		FnOut:       "key.txt",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	fs := memfs.NewMemFs()
	key := []byte("test")
	err := s.Age.EncryptItem(fs, &p, key)
	s.Require().Nil(err)

	encOut, err := fs.ReadFileOutPath(&p)
	s.Require().Nil(err)
	s.Require().NotEmpty(encOut)
	fsDec := memfs.NewMemFs()

	p.FnIn = "key.txt.age"
	err = fsDec.MakeFileIn(&p, encOut)
	s.Require().Nil(err)

	err = s.Age.DecryptToMemFsFile(&p, fsDec)
	s.Require().Nil(err)

	decOut, err := fsDec.ReadFileOutPath(&p)
	s.Require().Nil(err)
	s.Require().NotEmpty(decOut)
	s.Require().Equal(string(key), string(decOut))
	fmt.Println(string(decOut))
}

// use age-keygen -o private_key.txt to create a pubkey/private key pair for here
func (s *AgeEncryptionTestSuite) TestDecryption() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "kube.tar.gz.age",
		FnOut:       "kube_decrypted.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	err := s.Age.DecryptToFile(&p)
	s.Require().Nil(err)
}

func (s *AgeEncryptionTestSuite) TestDecryptionToInMemFs() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		FnIn:        "kube.tar.gz.age",
		FnOut:       "kube_decrypted.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	m := memfs.NewMemFs()
	err := s.Age.DecryptToMemFsFile(&p, m)
	s.Require().Nil(err)
}

func TestAgeEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(AgeEncryptionTestSuite))
}
