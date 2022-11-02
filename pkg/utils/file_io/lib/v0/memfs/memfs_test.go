package memfs

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MemFsTestSuite struct {
	suite.Suite
	MemFS MemFS
}

func (t *MemFsTestSuite) SetupTest() {
	t.MemFS = NewMemFs()
}

func (t *MemFsTestSuite) TestCreateFileMemFs() {

}

func (t *MemFsTestSuite) TestCreateMemFsAndStoreRead() {

}

func TestMemFsTestSuite(t *testing.T) {
	suite.Run(t, new(MemFsTestSuite))
}
