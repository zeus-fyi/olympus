package _struct

import (
	jen "github.com/dave/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/core/primitives"
)

func genDeclAt18(structToWrite primitives.StructGen) jen.Code {
	fields := make([]jen.Code, len(structToWrite.Fields))
	for i, f := range structToWrite.Fields {
		fields[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(structToWrite.Name).Struct(fields...)
	return _struct
}

func genFile(fw primitives.FileWrapper, structGen primitives.StructGen) *jen.File {
	ret := jen.NewFile(fw.FileName)
	ret.Add(genDeclAt18(structGen))
	return ret
}
