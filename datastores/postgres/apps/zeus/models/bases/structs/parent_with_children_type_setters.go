package structs

func (pc *SuperParentClass) SetChildClassTypeName(childClassTypeName string) {
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
