package _struct

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type StructTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/pkg/hera/cookbook/autogen/types_template_preview/structs"

func createTestCodeGenShell() lib.CodeGen {
	p := structs.Path{
		PackageName: "_struct",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "struct.go",
		Env:         "",
	}
	cg := lib.NewCodeGen(p)
	return cg
}

func (s *StructTestSuite) TestCodeGen() {
	cg := createTestCodeGenShell()
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

	cg.Add(genHeader())
	cg.Add(genDeclAt85())
	cg.Add(genDeclAt117())
	cg.Add(genDeclAt289())
	cg.Add(genFuncGetRowValues())
	cg.Add(structToMake.GenerateStructJenCode())
	err := cg.Save()
	s.Assert().Nil(err)

	s.Cleanup = false
	if s.Cleanup {
		s.DeleteFile(cg.Path.Fn)
	}
}

func TestFuncTestSuite(t *testing.T) {
	suite.Run(t, new(StructTestSuite))
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
