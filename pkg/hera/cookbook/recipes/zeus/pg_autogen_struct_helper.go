package zeus

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/file_shells/base"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
)

func genDeclAt18() jen.Code {
	return jen.Null()
}
func genDeclAt118() jen.Code {
	return jen.Null().Type().Id("Chart").Struct(jen.Qual("github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen", "ChartPackages"))
}

func GenerateZeusStruct(fw primitives.FileWrapper) error {
	f := base.FileBase(fw)
	f.Add(genDeclAt118())
	return f.Save(fw.PackageName)
}
