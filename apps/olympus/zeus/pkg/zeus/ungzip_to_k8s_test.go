package zeus

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type UnGzipToMemFsTestSuite struct {
	base.TestSuite
	memfs.MemFS
}

func (s *UnGzipToMemFsTestSuite) SetupTest() {
	s.MemFS = memfs.NewMemFs()
	s.InitLocalConfigs()
}

func (s *UnGzipToMemFsTestSuite) TestUnGzipIntoMemFs() {
	p := structs.Path{DirIn: "./", DirOut: "./", Fn: "zeus.tar.gz"}
	byteArray, err := ioutil.ReadFile(p.FileInPath())
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
