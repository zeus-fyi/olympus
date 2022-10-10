package base

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"

func structMock() primitives.StructGen {
	fieldOne := primitives.Field{
		Name: "ChartPackageID",
		Type: "int",
	}
	fieldTwo := primitives.Field{
		Name: "ChartName",
		Type: "string",
	}
	fieldThree := primitives.Field{
		Name: "ChartVersion",
		Type: "string",
	}

	structToMake := primitives.StructGen{
		Name:   "ChartPackages",
		Fields: nil,
	}
	structToMake.AddFields(fieldOne, fieldTwo, fieldThree)
	return structToMake
}
