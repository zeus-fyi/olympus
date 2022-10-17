package structs

import "github.com/zeus-fyi/jennifer/jen"

func GenCreateStructWithExternalStructInheritance(wrapperStructName, extStructPath, extStructName string) jen.Code {
	return jen.Null().Type().Id(wrapperStructName).Struct(jen.Qual(extStructPath, extStructName))
}
