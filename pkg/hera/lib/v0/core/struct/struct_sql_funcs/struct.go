package struct_sql_funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
)

type StructFuncGenSQL struct {
	primitives.StructGen
}

func (s *StructFuncGenSQL) GenerateGetRowValues(field ...primitives.Field) jen.Code {
	fnName := "GetRowValues"
	return genFuncGetRowValues(s.Name, fnName)
}

func genFuncGetRowValues(structName, fnName string) jen.Code {
	return jen.Func().Params(jen.Id("v").Op("*").Id(structName)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}
