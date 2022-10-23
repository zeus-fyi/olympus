package models

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

const Sn = "StructNameExample"

type StructNameExample struct {
	Field  string `json:"jsonName" db:"db_field_name" etc:"pattern"`
	FieldN int    `json:"fieldNameN" db:"db_field_name_N" etc:"patternN"`
}

type StructNameExamples []StructNameExample

func (v *StructNameExample) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	case "fieldGroup1":
		pgValues = apps.RowValues{v.Field}
	default:
		// should default to all
		pgValues = apps.RowValues{v.Field, v.FieldN}
	}
	return pgValues
}

func (v *StructNameExample) GetTableColumnValues(queryName string) []string {
	var columnValues []string
	switch queryName {
	case "fieldGroup1":
		columnValues = []string{"col1"}
	default:
		// should default to all
		columnValues = []string{"col1", "col2"}
	}
	return columnValues
}
