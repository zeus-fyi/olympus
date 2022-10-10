package primitives

type SwitchCase struct {
	SwitchName string
}

func genSwitchCase(caseVariableName string) {
	//jen.Switch(jen.Id("queryName")).Block(jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps")
	//.Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").
	//	Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN"))))
}
