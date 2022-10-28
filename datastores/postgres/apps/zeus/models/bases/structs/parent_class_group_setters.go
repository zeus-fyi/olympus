package structs

func (spg *SuperParentClassGroup) SetChartPackageResourceAndParentIDs(chartPackageID, chartResourceID, parentClassTypeID int) {
	spg.ChartSubcomponentParentClassTypes.ChartPackageID = chartPackageID
	spg.ChartSubcomponentParentClassTypes.ChartComponentResourceID = chartResourceID
	spg.SetSuperParentClassGroupParentID(parentClassTypeID)

	for i, sp := range spg.SuperParentClassSlice {
		sp.SetChartPackageAndResourceID(chartPackageID, chartResourceID)
		sp.SetParentClassTypeID(parentClassTypeID)
		spg.SuperParentClassSlice[i] = sp
	}

	for k, sp := range spg.SuperParentClassMap {
		sp.SetChartPackageAndResourceID(chartPackageID, chartResourceID)
		sp.SetParentClassTypeID(parentClassTypeID)
		spg.SuperParentClassMap[k] = sp
	}
}

func (spg *SuperParentClassGroup) SetSuperParentClassGroupParentID(id int) {
	spg.ChartSubcomponentParentClassTypeID = id
}
