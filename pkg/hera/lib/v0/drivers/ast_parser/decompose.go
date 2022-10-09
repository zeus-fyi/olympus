package ast_parser

import (
	"go/parser"
	"go/token"
)

func Decompose(b []byte) DeclKindCounter {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", b, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return AstParser(f)
}
