package compression

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type CompressionTestSuite struct {
	base.CoreTestSuite
	Comp Compression
}

func (c *CompressionTestSuite) SetupTest() {
	c.Comp = NewCompression()
}

func (c *CompressionTestSuite) TestTarGzip() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		Fn:          "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.CreateTarGzipArchiveDir(&p)
	c.Require().Nil(err)
}

func (c *CompressionTestSuite) TestUnGzip() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		Fn:          "./kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.UnGzip(&p)
	c.Require().Nil(err)
}

func (c *CompressionTestSuite) TestUnGzipInMemFS() {
	pkube := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		Fn:          "kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	m := memfs.NewMemFs()

	r := readers.ReaderLib{}
	b := r.ReadFile(pkube)
	c.Require().NotEmpty(b)

	err := m.MakeFile(&pkube, b)
	c.Require().Nil(err)

	ptgz := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		Fn:          "kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	err = c.Comp.UnGzipFromInMemFsOutToInMemFS(&ptgz, m)
	c.Require().Nil(err)
}

func (c *CompressionTestSuite) TestTar() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		Fn:          "kube.tar",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.TarFolder(&p)
	c.Require().Nil(err)
}

func TestCompressionTestSuite(t *testing.T) {
	suite.Run(t, new(CompressionTestSuite))
}
