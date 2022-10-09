package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
)

type Cookbook struct {
	lib.HeraLib
}

var (
	c            = Cookbook{}
	CookbookPath = structs.Path{
		DirIn:  "code_templates",
		DirOut: "autogen",
	}
	print = file_io.Printer{}
	log   = logging.Logger{}
)
