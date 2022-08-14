package _struct

import (
	jen "github.com/dave/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/codegen/cookbook/core/primitives"
)

func genDeclAt18(structToWrite primitives.StructGen) jen.Code {
	fields := make([]jen.Code, len(structToWrite.Fields))
	for i, f := range structToWrite.Fields {
		fields[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(structToWrite.Name).Struct(fields...)
	return _struct
}

func genFile(packageName string, structGen primitives.StructGen) *jen.File {
	ret := jen.NewFile(packageName)
	ret.Add(genDeclAt18(structGen))
	return ret
}
