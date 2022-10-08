package primitives

import "github.com/dave/jennifer/jen"

type FuncGen struct {
	Name         string
	Fields       []Field
	ReturnFields []Field
}

func (f *FuncGen) GetFieldStatement() *jen.Statement {
	statement := &jen.Statement{}
	for _, item := range f.Fields {
		statement.Add(jen.Id(item.Name).Id(item.Type))
	}
	return statement
}

func (f *FuncGen) AddField(field Field) {
	f.Fields = append(f.Fields, field)
}

func (f *FuncGen) AddFields(field []Field) {
	f.Fields = append(f.Fields, field...)
}

func (f *FuncGen) AddReturnField(field Field) {
	f.ReturnFields = append(f.Fields, field)
}

func (f *FuncGen) AddReturnFields(field []Field) {
	f.ReturnFields = append(f.Fields, field...)
}
