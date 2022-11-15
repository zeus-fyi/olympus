package memfs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type MemFsTestSuite struct {
	suite.Suite
	MemFS MemFS
}

func (t *MemFsTestSuite) SetupTest() {
	t.MemFS = NewMemFs()
}

func (t *MemFsTestSuite) TestCreateFileMemFs() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		FnIn:        "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := t.MemFS.MkPathDirAll(&p)
	t.Require().Nil(err)

}

func (t *MemFsTestSuite) TestCreateMemFsAndStoreRead() {

}

func TestMemFsTestSuite(t *testing.T) {
	suite.Run(t, new(MemFsTestSuite))
}
