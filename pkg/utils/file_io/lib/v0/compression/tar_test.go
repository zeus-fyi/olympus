package compression

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type TarTestSuite struct {
	CompressionTestSuite
}

func (c *Lz4TestSuite) TestTar() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		FnIn:        "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.TarCompress(&p)
	c.Require().Nil(err)
}

func (c *Lz4TestSuite) TestUnTar() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		FnIn:        "kube.tar",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.TarDecompress(&p)
	c.Require().Nil(err)
}
func TestTarTestSuite(t *testing.T) {
	suite.Run(t, new(Lz4TestSuite))
}
