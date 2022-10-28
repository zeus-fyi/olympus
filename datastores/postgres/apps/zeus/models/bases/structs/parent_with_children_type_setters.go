package structs

func (pc *SuperParentClass) SetBothChildClassTypeNames(childClassTypeName string) {
	if pc.ChildClassSingleValue != nil {
		pc.ChildClassSingleValue.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeName = childClassTypeName
	}
	if pc.ChildClassMultiValue != nil {
		pc.ChildClassMultiValue.ChartSubcomponentChildClassTypeName = childClassTypeName
	}
}

func (pc *SuperParentClass) SetParentClassTypeName(parentClassTypeName string) {
	pc.ChartSubcomponentParentClassTypeName = parentClassTypeName
}

func (pc *SuperParentClass) SetSingleChildClassIdTypeNameKeyAndValue(classTypeID int, childClassTypeName, k, v string) {
	if pc.ChildClassSingleValue == nil {
		nc := NewChildClassSingleValue(childClassTypeName)
		nc.SetChildClassTypeIDs(classTypeID)
		pc.ChildClassSingleValue = &nc
	}
	pc.ChildClassSingleValue.SetSingleChildClassIDTypeNameKeyAndValue(classTypeID, childClassTypeName, k, v)
}
