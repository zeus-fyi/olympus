package template_driver

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/drivers/ast_parser"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/tojen/gen"
)

var l = logging.Logger{}
var p = file_io.FileIO{}

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

func (t *TemplateDriverLib) CreateCustomZeusTemplate(recipePath, cookbookPath structs.Path) error {
	retBytes, err := gen.GenerateFileBytes(p.ReadFile(recipePath), recipePath.PackageName, false, false)
	if err != nil {
		return err
	}
	ast_parser.Decompose(retBytes)
	recipePath.LeftExtendDirOutPath(cookbookPath.DirOut)
	err = p.CreateFile(recipePath, retBytes)
	if err != nil {
		return err
	}
	return err
}
