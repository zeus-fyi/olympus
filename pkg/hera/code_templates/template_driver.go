package code_templates

import (
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/tojen/gen"
)

var l = logging.Logger{}
var p = printer.Printer{}

func CreateJenFile(pathIn, pathOut structs.Path) error {
	retBytes, err := gen.GenerateFileBytes(p.ReadFile(pathIn), pathIn.PackageName, false, true)
	if l.ErrHandler(err) != nil {
		return err
	}
	err = p.CreateFile(pathOut, retBytes)
	return l.ErrHandler(err)
}
