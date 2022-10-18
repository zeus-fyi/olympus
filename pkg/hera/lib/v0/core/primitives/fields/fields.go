package fields

import "github.com/zeus-fyi/jennifer/jen"

type FileWrapper struct {
	PackageName string
	FileName    string
}

type Field struct {
	Pkg   string
	Name  string
	Type  string
	Value string
}

type CaseField struct {
	Name string
	Type string
	Body []*jen.Statement
}

func NewCaseField(name, typeName string) CaseField {
	var body []*jen.Statement
	return CaseField{
		Name: name,
		Type: typeName,
		Body: body,
	}
}

func (c *CaseField) AddBodyStatement(js *jen.Statement) {
	c.Body = append(c.Body, js)
	return
}
