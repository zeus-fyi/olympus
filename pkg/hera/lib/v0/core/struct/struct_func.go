package _struct

import (
	"github.com/zeus-fyi/jennifer/jen"
)

func CreateStructPtrFunc(structName, shortHand string) jen.Code {
	return jen.Func().Params(jen.Id(shortHand).Op("*").Id(structName))
}
