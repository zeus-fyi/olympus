package _struct

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/template_test"
)

type StructTestSuite struct {
	template_test.TemplateTestSuite
}

func (s *StructTestSuite) TestCodeGen() {
	fw := primitives.FileWrapper{PackageName: "_struct", FileName: "struct_example.go"}

	structToMake := primitives.StructGen{
		Name:   "StructExample",
		Fields: nil,
	}
	fieldOne := primitives.Field{
		Name: "IntField",
		Type: "int",
	}
	structToMake.AddField(fieldOne)

	fieldTwo := primitives.Field{
		Name: "StringField",
		Type: "string",
	}
	structToMake.AddField(fieldTwo)

	resp := genMutateFile(fw, structToMake)
	s.Assert().NotEmpty(resp)

	err := resp.Save(fw.FileName)
	s.Assert().Nil(err)

	s.Cleanup = true
	if s.Cleanup {
		s.DeleteFile(fw.FileName)
	}
}

func genMutateFile(fw primitives.FileWrapper, structGen primitives.StructGen) *jen.File {
	f := jen.NewFile(fw.PackageName)
	ret := genTestBase(f)
	ret.Add(AddStruct(structGen))
	return ret
}

func AddStruct(structToWrite primitives.StructGen) jen.Code {
	fields := make([]jen.Code, len(structToWrite.Fields))
	for i, f := range structToWrite.Fields {
		fields[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(structToWrite.Name).Struct(fields...)
	return _struct
}

func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(StructTestSuite))
}

func genTestBase(ret *jen.File) *jen.File {
	ret.Add(genHeader())
	ret.Add(genDeclAt85())
	ret.Add(genDeclAt117())
	ret.Add(genDeclAt289())
	ret.Add(genFuncGetRowValues())
	return ret
}

func genHeader() jen.Code {
	return jen.Null()
}
func genDeclAt85() jen.Code {
	return jen.Null().Var().Id("Sn").Op("=").Lit("StructNameExample")
}
func genDeclAt117() jen.Code {
	return jen.Null().Type().Id("StructNameExample").Struct(jen.Id("Field").Id("string"), jen.Id("FieldN").Id("int"))
}
func genDeclAt289() jen.Code {
	return jen.Null().Type().Id("StructNameExamples").Index().Id("StructNameExample")
}
func genFuncGetRowValues() jen.Code {
	return jen.Func().Params(jen.Id("v").Op("*").Id("StructNameExample")).Id("GetRowValues").Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}
