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

func (spg *SuperParentClassGroup) SetParentClassTypeNames(parentClassTypeName string) {
	spg.ChartSubcomponentParentClassTypes.ChartSubcomponentParentClassTypeName = parentClassTypeName
	for i, _ := range spg.SuperParentClassSlice {
		spg.SuperParentClassSlice[i].ChartSubcomponentParentClassTypeName = parentClassTypeName
	}
	for k, sp := range spg.SuperParentClassMap {
		sp.SetParentClassTypeName(parentClassTypeName)
		spg.SuperParentClassMap[k] = sp
	}
}
