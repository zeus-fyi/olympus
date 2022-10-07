package structs

import (
	autogen_structs2 "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type Package struct {
	ChartComponentKindName   string
	ChartComponentApiVersion string
	ChartSubcomponents       PackageSubcomponent
}

type PackageSubcomponent struct {
	ChartSubcomponentParentClassTypeId   int
	ChartSubcomponentParentClassTypeName string
	ChartSubcomponentChildClassTypeName  string
	ChartSubcomponentChildClassTypeId    int
	ChartSubcomponentKeyName             *string
	ChartSubcomponentValue               *string
	ChartSubcomponentFieldName           *string
	ChartSubcomponentJsonbKeyValues      *string
}

type PackageComponentMap map[string]map[int][]PackageSubcomponent

type ChildValuesSlice []autogen_structs2.ChartSubcomponentsChildValues
type VolumesSlice []autogen_structs2.Volumes
