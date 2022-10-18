package conditionals

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type SwitchCase struct {
	SwitchName string
	Conditions map[string]fields.CaseField
}

func NewSwitchCase(name string) SwitchCase {
	s := SwitchCase{
		SwitchName: name,
		Conditions: make(map[string]fields.CaseField),
	}
	return s
}

func (s *SwitchCase) AddCondition(cf fields.CaseField) {
	if _, ok := s.Conditions[cf.Name]; !ok {
		s.Conditions[cf.Name] = cf
	}
	s.Conditions[cf.Name] = cf
	return
}

func (s *SwitchCase) GenerateSwitchStatement() *jen.Statement {
	switchStatement := jen.Switch(jen.Id(s.SwitchName)).Block(s.generateSwitchBodyCaseStatement())
	return switchStatement
}
