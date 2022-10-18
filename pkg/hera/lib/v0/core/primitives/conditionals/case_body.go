package conditionals

import "github.com/zeus-fyi/jennifer/jen"

func (s *SwitchCase) generateSwitchBodyCaseStatement() *jen.Statement {
	caseStatements := jen.Null()
	for _, cond := range s.Conditions {
		casePrefix := s.generateCasePrefix(cond.Name)
		for _, body := range cond.Body {
			caseStatements.Add(casePrefix.Block(body))
		}
	}
	return caseStatements
}

func (s *SwitchCase) generateCasePrefix(condName string) *jen.Statement {
	caseStatement := jen.Case()
	if condName == "default" {
		caseStatement = jen.Default()
		return caseStatement
	}
	return caseStatement.Add(jen.Lit(condName))
}
