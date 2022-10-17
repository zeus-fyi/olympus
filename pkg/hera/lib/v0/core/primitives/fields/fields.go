package fields

import "github.com/zeus-fyi/jennifer/jen"

type FileWrapper struct {
	PackageName string
	FileName    string
}

type Field struct {
	Name string
	Type string
}

type CaseField struct {
	Name string
	Type string
	Body *jen.Statement
}
