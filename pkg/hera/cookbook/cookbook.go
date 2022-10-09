package cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
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
	print = printer.Printer{}
	log   = logging.Logger{}
)
