package struct_sql_funcs

//
//type StructFuncGenSQL struct {
//}
//
//func (s *StructFuncGenSQL) GenerateGetRowValues(name string, field ...fields.Field) jen.Code {
//	fnName := "GetRowValues"
//	return genFuncGetRowValues(name, fnName)
//}
//
//func genFuncGetRowValues(structName, fnName string) jen.Code {
//	return jen.Func().Params(jen.Id("v").Op("*").Id(structName)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
//}

//return jen.Func().Params(jen.Id("v").Op("*").Id(structName)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).
//	Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).
//		Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))),
//		jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"),
//			jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
