package _struct

import (
	"github.com/zeus-fyi/jennifer/jen"
)

func CreateStructVariableDeclAssignment(structName string) jen.Code {
	return jen.Null().Var().Id("Sn").Op("=").Lit(structName)
}

func CreateStructSliceVariableDecl(structName string) jen.Code {
	return jen.Null().Type().Id(structName).Index().Id(structName)
}

func CreateStructPtrFunc(structName, shortHand string) jen.Code {
	return jen.Func().Params(jen.Id(shortHand).Op("*").Id(structName))
}
