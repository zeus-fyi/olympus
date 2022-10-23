package vars

import "github.com/zeus-fyi/jennifer/jen"

func CreateStringAssignment(varName, litValue string) *jen.Statement {
	return jen.Id(varName).Op("=").Lit(litValue)
}
