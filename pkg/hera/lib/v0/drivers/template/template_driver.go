package template_driver

import (
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/tojen/gen"
)

var l = logging.Logger{}
var p = printer.Printer{}

type TemplateDriverLib struct {
}

func (t *TemplateDriverLib) CreateTemplate(recipePath, cookbookPath structs.Path) error {
	retBytes, err := gen.GenerateFileBytes(p.ReadFile(recipePath), recipePath.PackageName, false, false)
	if err != nil {
		return err
	}
	recipePath.LeftExtendDirOutPath(cookbookPath.DirOut)
	err = p.CreateFile(recipePath, retBytes)
	if err != nil {
		return err
	}
	return err
}
