package template

import (
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/tojen/gen"
)

var l = logging.Logger{}
var p = printer.Printer{}

func CreateJenFile(path structs.Path) error {
	retBytes, err := gen.GenerateFileBytes(p.ReadFile(path), path.PackageName, false, true)
	if l.ErrHandler(err) != nil {
		return err
	}
	f := gen.GenerateFile(retBytes, path.PackageName, false)
	err = f.Save(path.FileOutPath())
	return l.ErrHandler(err)
}

func CreateJenFilesFromDir(path structs.Path) error {
	pathsIn := p.BuildPathsFromDirInPath(path, ".go")
	for _, pathIn := range pathsIn.Slice {
		err := CreateJenFile(pathIn)
		if l.ErrHandler(err) != nil {
			return err
		}
	}
	return nil
}
