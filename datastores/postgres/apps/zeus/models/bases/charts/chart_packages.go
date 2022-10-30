package charts

import (
	"database/sql"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type Chart struct {
	autogen_bases.ChartPackages
	autogen_bases.ChartComponentResources
}

func NewChart() Chart {
	c := Chart{autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        "",
		ChartVersion:     "",
		ChartDescription: sql.NullString{},
	}, autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 0,
		ChartComponentKindName:   "",
		ChartComponentApiVersion: "",
	}}
	return c
}

const Sn = "Chart"

func (c *Chart) GetChartPackageID() int {
	return c.ChartPackageID
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
