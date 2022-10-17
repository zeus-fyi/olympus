package conditionals

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type SwitchCase struct {
	SwitchName string
	Conditions []fields.CaseField
}

func (s *SwitchCase) GenerateSwitchStatement() jen.Code {
	switchStatement := jen.Switch(jen.Id(s.SwitchName)).Block(s.generateSwitchBodyCaseStatement())
	return switchStatement
}

func (s *SwitchCase) generateSwitchBodyCaseStatement() *jen.Statement {
	//caseStatements := jen.Case(jen.Lit("fieldGroup1")).Block(jen.Id("pgValues").Op("=").Id("apps")
	//.Dot("RowValues").Values(jen.Id("v").Dot("Field"))), jen.Default().Block(jen.Id("pgValues").Op("=").Id("apps").
	//	Dot("RowValues").Values(jen.Id("v").Dot("Field"), jen.Id("v").Dot("FieldN"))))
	caseStatements := jen.Case()
	for _, cond := range s.Conditions {
		caseStatements.Add(jen.Lit(cond.Name).Block(cond.Body))
	}
	return caseStatements
}
