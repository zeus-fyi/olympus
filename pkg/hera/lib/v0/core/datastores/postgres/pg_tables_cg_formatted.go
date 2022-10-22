package postgres

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/table_formatting"
)

func (d *PgSchemaAutogen) ConvertTablesToCodeGenStructs() {
	tables := d.TableMap
	for _, tbl := range tables {
		tblName := table_formatting.FormatTableName(d.Settings, tbl)
		s := primitive.StructGen{
			Name:       tblName,
			Fields:     nil,
			PluralDecl: nil,
		}
		fieldsToAdd := make([]fields.Field, len(tbl.Columns))
		for i, col := range tbl.Columns {
			fieldsToAdd[i] = d.processTableElement(tbl, col)
		}
		s.AddFields(fieldsToAdd...)
		d.StructMapToCodeGen[tbl.Name] = s
	}
}

func (d *PgSchemaAutogen) processTableElement(tbl *database.Table, col database.Column) fields.Field {
	goType, _ := table_formatting.MapDbColumnTypeToGoType(d.Settings, d.Postgresql, col)
	colName, err := table_formatting.FormatColumnName(d.Settings, col.Name, tbl.Name)
	if err != nil {
		panic(err)
	}
	field := fields.NewFieldFromDB(tbl, col, colName, goType)
	return field
}
