package sql_query

import "github.com/zeus-fyi/jennifer/jen"

// GenPGExecSQLStatement this template relies on a common Pg package and it's var naming used across apps
func genPGExecSQLStatement() *jen.Statement {
	return jen.Id("err").Op(":=").Id("apps").Dot("Pg").Dot("Exec").Call(jen.Id("ctx"), jen.Id("q").Dot("SelectQuery").Call())
}

// genPGSqlExecRowsAffectedHandler this template relies on common var naming syntax used across apps
func genPGSqlExecRowsAffectedHandler() *jen.Statement {
	return jen.Id("rowsAffected").Op(":=").Id("r").Dot("RowsAffected").Call()
}

// TODO needs struct name and parent pkg name
func genPGSqlExecRowsAffectedDebugLog() *jen.Statement {
	return jen.Id("log").Dot("Debug").Call().Dot("Msgf").Call(jen.Lit("StructNameExamples: %s, Rows Affected: %d"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")), jen.Id("rowsAffected"))
}

func GenPGGenericExec() []jen.Code {
	return []jen.Code{
		genLogHeader(),
		genPGExecSQLStatement(),
		genSqlExecErrHandler(),
		genPGSqlExecRowsAffectedHandler(),
		genPGSqlExecRowsAffectedDebugLog(),
		genSqlInsertReturnWrappedErr(),
	}
}

func genLogHeader() *jen.Statement {
	return jen.Id("log").Dot("Debug").Call().Dot("Interface").Call(jen.Lit("InsertQuery:"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}

func genSqlExecErrHandler() *jen.Statement {
	return jen.If(jen.Id("returnErr").Op(":=").Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").
		Dot("LogHeader").Call(jen.Id("models").
		Dot("Sn"))), jen.Id("returnErr").Op("!=").Id("nil").Block(jen.Return().Id("err")))
}

func genSqlInsertReturnWrappedErr() *jen.Statement {
	return jen.Return().Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}
