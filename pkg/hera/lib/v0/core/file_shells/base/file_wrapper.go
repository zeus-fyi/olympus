package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
)

func (b *FileComponentBaseElements) GenerateFileShell(fw primitives.FileWrapper) *jen.File {
	f := jen.NewFile(fw.PackageName)
	f.Add(genHeader())
	return f
}

func genHeader() jen.Code {
	return jen.Null()
}
