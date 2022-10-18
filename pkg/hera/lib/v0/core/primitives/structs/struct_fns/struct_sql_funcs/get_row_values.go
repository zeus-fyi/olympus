package struct_sql_funcs

type StructFuncGenSQL struct {
}

//func (s *StructFuncGenSQL) GenerateGetRowValues(structGen structs.StructGen) jen.Code {
//	fnName := "GetRowValues"
//	return genFuncGetRowValues(structGen, fnName)
//}

//func genFuncGetRowValues(structGen structs.StructGen, fnName string) jen.Code {
//	v := vars.VariableGen{}
//	v.GenStructInstructs[structGen.Name] = structGen
//	declStruct := v.CreateStructDecl("pgValues", "apps", structGen, true)
//
//	return jen.Func().Params(jen.Id("v").Op("*").Id(structGen.Name)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).Block(declStruct, jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
//}

//return jen.Func().Params(jen.Id("v").Op("*").Id(structName)).Id(fnName).Params(jen.Id("queryName").Id("string")).Params(jen.Id("apps").Dot("RowValues")).
//	Block(jen.Id("pgValues").Op(":=").Id("apps").Dot("RowValues").Values(), jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).
//		Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"))),
//		jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").Dot("RowValues").Values(jen.Id("v").Dot("Field"),
//			jen.Id("v").Dot("FieldN")))), jen.Return().Id("pgValues"))
