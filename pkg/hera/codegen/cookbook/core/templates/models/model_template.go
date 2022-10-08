package models

import "github.com/zeus-fyi/olympus/datastores/postgres_apps"

const sn = "StructNameExample"

type StructNameExample struct {
	Field  string `json:"jsonName" db:"db_field_name" etc:"pattern"`
	FieldN int    `json:"fieldNameN" db:"db_field_name_N" etc:"patternN"`
}

type StructNameExamples []StructNameExample

func (v *StructNameExample) GetRowValues(queryName string) postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{}
	switch queryName {
	case "fieldGroup1":
		pgValues = postgres_apps.RowValues{v.Field}
	default:
		// should default to all
		pgValues = postgres_apps.RowValues{v.Field, v.FieldN}
	}
	return pgValues
}
