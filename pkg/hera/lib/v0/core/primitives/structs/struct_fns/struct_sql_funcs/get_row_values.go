package struct_sql_funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs/struct_fns"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

type StructFuncGenSQL struct {
	struct_fns.StructFn
}

func (s *StructFuncGenSQL) GenerateGetRowValues(structGen structs.StructGen, field ...fields.Field) jen.Code {
	fnName := "GetRowValues"
	return genFuncGetRowValues(structGen, fnName)
}

func genFuncGetRowValues(structGen structs.StructGen, fnName string) jen.Code {
	v := vars.VariableGen{}
	v.Structs[structGen.Name] = structGen
	declStruct := v.CreateStructDecl("pgValues", "apps", structGen.Name, true)
	return jen.Func().Params(jen.Id("v").Op("*").Id(structGen.Name)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(declStruct, jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}

//return jen.Func().Params(jen.Id("v").Op("*").Id(structName)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).
//	Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).
//		Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))),
//		jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"),
//			jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
