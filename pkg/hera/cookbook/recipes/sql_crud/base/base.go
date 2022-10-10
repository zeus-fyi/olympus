package base

import (
	"github.com/zeus-fyi/jennifer/jen"
)

func genHeader() jen.Code {
	return jen.Null()
}
func genDeclAt85(structName string) jen.Code {
	return jen.Null().Var().Id("Sn").Op("=").Lit(structName)
}

func genFuncGetRowValues() jen.Code {
	return jen.Func().Params(jen.Id("v").Op("*").Id("StructNameExample")).Id("GetRowValues").Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
}
