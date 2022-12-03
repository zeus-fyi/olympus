package compression

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Lz4TestSuite struct {
	CompressionTestSuite
}

func (c *Lz4TestSuite) TestLz4Cmp() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./in",
		DirOut:      "./cmp",
		FnIn:        "gfdamnit.json",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	r := readers.ReaderLib{}

	err := c.Comp.Lz4CompressDir(&p)
	c.Require().Nil(err)
	p.FnOut = "gfdamnit.json.tar.lz4"
	//bc := r.ReadFilePathOutJsonObject(p)

	p = filepaths.Path{
		PackageName: "",
		DirIn:       "./in",
		DirOut:      "./out",
		FnIn:        "gfdamnit.json",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	inMem := memfs.NewMemFs()
	b := r.ReadJsonObject(p)

	err = inMem.MakeFileIn(&p, b)
	c.Require().Nil(err)

	fs, err := c.Comp.Lz4CompressInMemFsFile(&p, inMem)
	c.Require().Nil(err)
	c.Require().NotEmpty(fs)

	out, err := fs.ReadFileOutPath(&p)
	c.Require().Nil(err)

	err = os.WriteFile(p.FileOutPath(), out, 0644)
	c.Require().Nil(err)

	p.FnOut = "gfdamnit.json"
	inMem, err = c.Comp.Lz4DecompressInMemFsFile(&p, inMem)
	c.Require().Nil(err)

	decOut, err := fs.ReadFileOutPath(&p)
	c.Require().Nil(err)
	err = os.WriteFile(p.FileOutPath(), decOut, 0644)
	c.Require().Nil(err)
}

func (c *Lz4TestSuite) TestLz4Dec() {
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./cmp",
		DirOut:      "./out",
		FnIn:        "validator-balance-epoch-163999.json.tar.lz4",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	inMem := memfs.NewMemFs()
	r := readers.ReaderLib{}

	b := r.ReadJsonObject(p)

	err := inMem.MakeFileIn(&p, b)
	c.Require().Nil(err)
	inMem, err = c.Comp.Lz4DecompressInMemFsFile(&p, inMem)
	//err := c.Comp.Lz4Decompress(&p)
	c.Require().Nil(err)
	c.Assert().NotEmpty(inMem)

}

func TestLz4TestSuite(t *testing.T) {
	suite.Run(t, new(Lz4TestSuite))
}
