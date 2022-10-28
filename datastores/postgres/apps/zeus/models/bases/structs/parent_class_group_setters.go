package structs

func (sps *SuperParentClassGroup) SetChartPackageAndResourceID(chartPackageID, chartResourceID int) {
	sps.ChartSubcomponentParentClassTypes.ChartPackageID = chartPackageID
	sps.ChartSubcomponentParentClassTypes.ChartComponentResourceID = chartResourceID

	for _, sp := range sps.SuperParentClassSlice {
		sp.SetChartPackageAndResourceID(chartPackageID, chartResourceID)
	}

	for k, sp := range sps.SuperParentClassMap {
		sp.SetChartPackageAndResourceID(chartPackageID, chartResourceID)
		sps.SuperParentClassMap[k] = sp
	}
}

func (sps *SuperParentClassGroup) SetParentClassTypeNames(parentClassTypeName string) {
	sps.ChartSubcomponentParentClassTypes.ChartSubcomponentParentClassTypeName = parentClassTypeName
	for i, _ := range sps.SuperParentClassSlice {
		sps.SuperParentClassSlice[i].ChartSubcomponentParentClassTypeName = parentClassTypeName
	}
	for k, sp := range sps.SuperParentClassMap {
		sp.SetParentClassTypeName(parentClassTypeName)
		sps.SuperParentClassMap[k] = sp
	}
}

func (sps *SuperParentClassGroup) SetParentClassTypeIDs(parentClassTypeID int) {
	sps.ChartSubcomponentParentClassTypes.ChartSubcomponentParentClassTypeID = parentClassTypeID

	for _, sp := range sps.SuperParentClassSlice {
		sp.SetParentClassTypeID(parentClassTypeID)
	}
	for k, sp := range sps.SuperParentClassMap {
		sp.SetParentClassTypeID(parentClassTypeID)
		sps.SuperParentClassMap[k] = sp
	}
}
