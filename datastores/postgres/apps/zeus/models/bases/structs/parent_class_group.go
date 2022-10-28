package structs

type SuperParentClassGroup struct {
	SuperParentClassSlice
	SuperParentClassMap
}

type SuperParentClassSlice []SuperParentClass
type SuperParentClassMap map[int]SuperParentClass

func NewSuperParentClassGroup() SuperParentClassGroup {
	spcg := SuperParentClassGroup{
		SuperParentClassSlice: []SuperParentClass{},
		SuperParentClassMap:   make(map[int]SuperParentClass),
	}
	return spcg
}
