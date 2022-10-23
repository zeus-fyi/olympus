package postgres

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/table_formatting"
)

func (d *PgSchemaAutogen) ConvertTablesToCodeGenStructs() {
	tables := d.TableMap
	for _, tbl := range tables {
		if string_utils.FilterStringWithOpts(tbl.Name, d.Filter) {
			tblName := table_formatting.FormatTableName(d.Settings, tbl)
			s := primitive.StructGen{
				Name:       tblName,
				Fields:     nil,
				PluralDecl: nil,
				DBTable:    tbl,
			}
			fieldsToAdd := make(map[string]fields.Field)
			for _, col := range tbl.Columns {
				fieldsToAdd[col.Name] = d.processTableElement(tbl, col)
			}
			s.AddFieldsFromMap(fieldsToAdd)
			d.StructMapToCodeGen[tbl.Name] = s
		}
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
