package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (b *FileComponentBaseElements) GenerateFileShell(path filepaths.Path) *jen.File {
	f := jen.NewFile(path.PackageName)
	return f
}
