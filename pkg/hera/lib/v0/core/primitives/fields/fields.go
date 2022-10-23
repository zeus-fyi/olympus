package fields

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
)

type FileWrapper struct {
	PackageName string
	FileName    string
}

type Field struct {
	Pkg   string
	Name  string
	Type  string
	Value string

	DbMetadata
	FnField *jen.Statement
}

func (f *Field) DbFieldName() string {
	return f.DbMetadata.Column.Name
}

func (f *Field) DbTableName() string {
	return f.DbMetadata.Table.Name
}

func NewFieldFromDB(tbl *database.Table, col database.Column, colName, goType string) Field {
	newDbField := NewDbMetadata(tbl, col)
	f := Field{
		Pkg:        "",
		Name:       colName,
		Type:       goType,
		Value:      "",
		DbMetadata: newDbField,
		FnField:    nil,
	}
	return f
}
