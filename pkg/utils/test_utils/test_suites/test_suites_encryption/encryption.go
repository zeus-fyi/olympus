package test_suites_encryption

import (
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type EncryptionTestSuite struct {
	test_suites_base.TestSuite

	Age  encryption.Age
	Comp compression.Compression
}

func (s *EncryptionTestSuite) SetupTest() {
	s.SetupLocalAge()
	s.Comp = compression.NewCompression()
}

func (s *EncryptionTestSuite) SetupLocalAge() {
	s.Tc = configs.InitLocalTestConfigs()
	pubKey := s.Tc.LocalAgePubkey
	privKey := s.Tc.LocalAgePkey
	s.Age = encryption.NewAge(privKey, pubKey)
}
