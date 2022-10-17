package base

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

func structMock() structs.StructGen {
	fieldOne := fields.Field{
		Name: "ChartPackageID",
		Type: "int",
	}
	fieldTwo := fields.Field{
		Name: "ChartName",
		Type: "string",
	}
	fieldThree := fields.Field{
		Name: "ChartVersion",
		Type: "string",
	}

	structToMake := structs.StructGen{
		Name:   "ChartPackage",
		Fields: nil,
	}
	structToMake.AddFields(fieldOne, fieldTwo, fieldThree)
	return structToMake
}
