package charts

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type Chart struct {
	autogen_bases.ChartPackages
}

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
