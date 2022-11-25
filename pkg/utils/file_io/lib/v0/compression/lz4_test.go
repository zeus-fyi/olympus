package compression

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Lz4TestSuite struct {
	CompressionTestSuite
}

func (c *Lz4TestSuite) TestLz4Cmp() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		FnIn:        "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.Lz4CompressDir(&p)
	c.Require().Nil(err)
}

func (c *Lz4TestSuite) TestLz4Dec() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		FnIn:        "kube.tar.lz4",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.Lz4Decompress(&p)
	c.Require().Nil(err)
}

func TestLz4TestSuite(t *testing.T) {
	suite.Run(t, new(Lz4TestSuite))
}
