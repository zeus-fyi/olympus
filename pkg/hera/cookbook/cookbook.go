package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
)

type Cookbook struct {
	lib.HeraLib
}

var (
	c            = Cookbook{}
	CookbookPath = filepaths.Path{
		DirIn:  "code_templates",
		DirOut: "autogen",
	}
	fileIO = file_io.FileIO{}
	log    = logging.Logger{}
)
