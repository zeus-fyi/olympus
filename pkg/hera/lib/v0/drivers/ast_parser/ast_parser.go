package ast_parser

import (
	"go/ast"
)

type DeclKindCounter map[string]int

func AstParser(f *ast.File) DeclKindCounter {
	// Print the location and kind of each declaration in f.
	m := make(DeclKindCounter)
	for _, decl := range f.Decls {
		ProcessDeclTypes(decl, m)
	}
	return m
}

func ProcessDeclTypes(d ast.Decl, m DeclKindCounter) DeclKindCounter {
	switch d.(type) {
	case *ast.GenDecl:
		m = GenDeclTypeHandler(d.(*ast.GenDecl), m)
	case *ast.FuncDecl:
		m = GenFuncDeclTypeHandler(d.(*ast.FuncDecl), m)
	}
	return m
}

func GenDeclTypeHandler(gen *ast.GenDecl, m DeclKindCounter) DeclKindCounter {
	kind := gen.Tok.String()
	AddCountToKey(kind, m)
	return m
}

func GenFuncDeclTypeHandler(fun *ast.FuncDecl, m DeclKindCounter) DeclKindCounter {
	kind := fun.Name.String()
	AddCountToKey(kind, m)
	return m
}

func AddCountToKey(key string, m DeclKindCounter) DeclKindCounter {

	if c, okay := m[key]; !okay {
		m[key] = 1
	} else {
		m[key] = c + 1
	}
	return m
}
