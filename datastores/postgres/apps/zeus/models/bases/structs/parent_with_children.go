package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type SuperParentClass struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	*ChildClassSingleValue
	*ChildClassMultiValue
}

func NewSuperParentClass(parentClassTypeName string) SuperParentClass {
	sc := SuperParentClass{
		ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
			ChartPackageID:                       0,
			ChartComponentResourceID:             0,
			ChartSubcomponentParentClassTypeID:   0,
			ChartSubcomponentParentClassTypeName: parentClassTypeName,
		},
		ChildClassSingleValue: nil,
		ChildClassMultiValue:  nil,
	}
	return sc
}

func NewSuperParentClassWithMultiChildType(parentClassTypeName, multiChildClassTypeName string) SuperParentClass {
	ccmv := NewChildClassMultiValues(multiChildClassTypeName)
	sc := SuperParentClass{
		ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
			ChartPackageID:                       0,
			ChartComponentResourceID:             0,
			ChartSubcomponentParentClassTypeID:   0,
			ChartSubcomponentParentClassTypeName: parentClassTypeName,
		},
		ChildClassMultiValue: &ccmv,
	}
	return sc
}

func NewSuperParentClassWithBothChildTypes(parentClassTypeName, singleChildClassType, multiChildClassTypeName string) SuperParentClass {
	ccmv := NewChildClassMultiValues(multiChildClassTypeName)
	ccsv := NewChildClassSingleValue(singleChildClassType)
	sc := SuperParentClass{
		ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
			ChartPackageID:                       0,
			ChartComponentResourceID:             0,
			ChartSubcomponentParentClassTypeID:   0,
			ChartSubcomponentParentClassTypeName: parentClassTypeName,
		},
		ChildClassSingleValue: &ccsv,
		ChildClassMultiValue:  &ccmv,
	}
	return sc
}
