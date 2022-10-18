package common

import "github.com/zeus-fyi/jennifer/jen"

func GenPGExecSQLStatementErrHandler() *jen.Statement {
	return jen.If(jen.Id("returnErr").Op(":=").Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").
		Dot("LogHeader").Call(jen.Id("models").
		Dot("Sn"))), jen.Id("returnErr").Op("!=").Id("nil").Block(jen.Return().Id("err")))
}
