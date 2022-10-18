package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/conditionals"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
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

func genFuncGetRowValues2(structGen primitive.StructGen, fnName string) jen.Code {
	v := vars.NewVarGen()
	v.InsertStruct(genRowValuesStructTemplate())
	declInitStruct := DeclarePgValuesStructVar(v, "init", v.GenStructInstructs["init"])
	declDefaultCaseFields := DeclarePgValuesStructVar(v, "embedded", structGen)

	sc := conditionals.NewSwitchCase("queryName")
	cf := fields.NewCaseField("default", "")
	cf.AddBodyStatement(declDefaultCaseFields)
	sc.AddCondition(cf)
	statement := sc.GenerateSwitchStatement()
	fn := funcs.NewFn(fnName)
	fn.Fields = []fields.Field{fields.Field{
		Name:  "queryName",
		Type:  "string",
		Value: "",
	}}
	fn.ReturnFields = []fields.Field{fields.Field{
		Pkg:   "apps",
		Name:  "pgValues",
		Type:  "RowValues",
		Value: "",
	}}

	fn.AddBodyStatement(declInitStruct)
	fn.AddBodyStatement(statement)

	return fn.GenerateStructPtrFunc(structGen)

}

//return jen.Func().Params(jen.Id(structGen.ShortHand()).Op("*").Id(structGen.Name)).Id(fnName).
//Params(jen.Id("queryName").Id("string")).
//Params(jen.Id("apps").Dot("RowValues")).
//Block(declInitStruct,
//jen.Switch(jen.Id("queryName")).
//Block(jen.Case(jen.Lit("fieldGroup1")).
//Block(declDefaultCaseFields),
//jen.Default().Block(declDefaultCaseFields)),
//jen.Return().Id("pgValues"))

func genFuncGetRowValues() jen.Code {
	return jen.Func().Params(jen.Id("v").Op("*").Id("StructNameExample")).Id("GetRowValues").
		Params(jen.Id("queryName").Id("string")).
		Params(jen.Id("apps").Dot("RowValues")).
		Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(),
			jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).
				Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))),
				jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"),
					jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}
