package structs

func (pc *SuperParentClass) GetSuperParentClassTypeID() int {
	return pc.ChartSubcomponentParentClassTypeID
}

func (pc *SuperParentClass) SetParentClassTypeID(id int) {
	pc.ChartSubcomponentParentClassTypeID = id
	if pc.ChildClassSingleValue != nil {
		pc.ChildClassSingleValue.ChartSubcomponentParentClassTypeID = id
	}
	if pc.ChildClassMultiValue != nil {
		pc.ChildClassMultiValue.ChartSubcomponentParentClassTypeID = id
	}
}

func (pc *SuperParentClass) SetBothChildClassTypeID(id int) {
	if pc.ChildClassSingleValue != nil {
		pc.ChildClassSingleValue.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id
	}
	if pc.ChildClassMultiValue != nil {
		if len(pc.ChildClassMultiValue.Values) > 0 {
			for i, _ := range pc.ChildClassMultiValue.Values {
				pc.ChildClassMultiValue.Values[i].ChartSubcomponentChildClassTypeID = id
			}
		}
	}
}

func (pc *SuperParentClass) SetChartPackageAndResourceID(chartPackageID, chartComponentResourceID int) {
	pc.ChartPackageID = chartPackageID
	pc.ChartComponentResourceID = chartComponentResourceID
}
