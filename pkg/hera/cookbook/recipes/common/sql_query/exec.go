package sql_query

import "github.com/zeus-fyi/jennifer/jen"

// GenPGExecSQLStatement this template relies on a common Pg package and it's var naming used across apps
func GenPGExecSQLStatement() *jen.Statement {
	return jen.Id("err").Op(":=").Id("apps").Dot("Pg").Dot("Exec").Call(jen.Id("ctx"), jen.Id("q").Dot("SelectQuery").Call())
}

// GenPGSqlExecRowsAffectedHandler this template relies on common var naming syntax used across apps
func GenPGSqlExecRowsAffectedHandler() *jen.Statement {
	return jen.Id("rowsAffected").Op(":=").Id("r").Dot("RowsAffected").Call()
}
