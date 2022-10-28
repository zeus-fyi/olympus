package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type SuperParentClassGroup struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	SuperParentClassSlice
	SuperParentClassMap
}

type SuperParentClassSlice []SuperParentClass
type SuperParentClassMap map[string]SuperParentClass

func NewSuperParentClassGroup(parentClassTypeName string) SuperParentClassGroup {
	spcg := SuperParentClassGroup{
		SuperParentClassSlice: []SuperParentClass{},
		SuperParentClassMap:   make(map[string]SuperParentClass),
	}
	spcg.SetParentClassTypeNames(parentClassTypeName)
	return spcg
}
