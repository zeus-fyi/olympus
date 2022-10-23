package vars

import "github.com/zeus-fyi/jennifer/jen"

func (v *VariableGen) CreateConstStringDecl(name string) *jen.Statement {
	if value, ok := v.StringConstants[name]; ok {
		return jen.Null().Const().Id(name).Op("=").Lit(value)
	}
	return jen.Null()
}
