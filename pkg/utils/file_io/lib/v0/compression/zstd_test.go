package compression

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type ZstdTestSuite struct {
	CompressionTestSuite
}

func (c *ZstdTestSuite) TestZstdCmp() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./.kube",
		DirOut:      "./",
		FnIn:        "kube",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.ZstCompressDir(&p)
	c.Require().Nil(err)
}

func (c *ZstdTestSuite) TestZstdDec() {
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./kube",
		FnIn:        "kube.tar.zst",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := c.Comp.ZstdDecompress(&p)
	c.Require().Nil(err)
}

func TestZstdTestSuite(t *testing.T) {
	suite.Run(t, new(ZstdTestSuite))
}
