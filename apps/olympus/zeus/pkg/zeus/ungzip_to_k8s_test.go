package zeus

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type UnGzipToMemFsTestSuite struct {
	test_suites_base.TestSuite
	memfs.MemFS
}

func (s *UnGzipToMemFsTestSuite) SetupTest() {
	s.MemFS = memfs.NewMemFs()
	s.InitLocalConfigs()
}

func (s *UnGzipToMemFsTestSuite) TestUnGzipIntoMemFs() {
	p := filepaths.Path{DirIn: "./", DirOut: "./", FnIn: "zeus.tar.gz"}
	byteArray, err := os.ReadFile(p.FileDirOutFnInPath())
	s.Require().Nil(err)
	b := &bytes.Buffer{}
	_, err = b.Write(byteArray)
	s.Require().Nil(err)
	nk, err := UnGzipK8sChart(b)
	s.Require().Nil(err)
	s.Assert().NotEmpty(nk)
}

func TestUnGzipToMemFsTestSuite(t *testing.T) {
	suite.Run(t, new(UnGzipToMemFsTestSuite))
}
