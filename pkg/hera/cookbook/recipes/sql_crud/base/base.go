package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelTemplate struct {
	lib.CodeGen
}

func NewModelTemplate(p structs.Path) ModelTemplate {
	m := ModelTemplate{lib.NewCodeGen(p)}
	return m
}

func (m *ModelTemplate) CreateTemplate() error {
	m.Add(m.Vars.CreateConstStringDecl(m.Path.PackageName))
	structJenCode := m.Structs.GenerateStructsJenCode(true)
	m.AddSlice(structJenCode)
	m.Add(genFuncGetRowValues())
	err := m.Save()
	return err
}

func (m *ModelTemplate) createRowValuesPtrFunc() error {
	m.Add(m.Vars.CreateConstStringDecl(m.Path.PackageName))
	structJenCode := m.Structs.GenerateStructsJenCode(true)
	m.AddSlice(structJenCode)
	m.Add(genFuncGetRowValues())
	err := m.Save()
	return err
}

func genFuncGetRowValues() jen.Code {
	return jen.Func().Params(jen.Id("v").Op("*").Id("StructNameExample")).Id("GetRowValues").Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}
