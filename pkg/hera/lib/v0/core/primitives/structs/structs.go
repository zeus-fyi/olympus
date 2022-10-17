package structs

import "github.com/zeus-fyi/jennifer/jen"

type StructsGen struct {
	Structs    []StructGen
	StructsMap map[string]StructGen
}

func (sg *StructsGen) AddStruct(s StructGen) {
	if len(sg.Structs) == 0 {
		sg.Structs = []StructGen{}
		sg.StructsMap = make(map[string]StructGen)
	}
	sg.StructsMap[s.Name] = s
	sg.Structs = append(sg.Structs, s)
}

func (sg *StructsGen) GenerateStructsJenCode(withPlural bool) []jen.Code {
	var structsCode []jen.Code
	for _, s := range sg.Structs {
		structsCode = append(structsCode, s.GenerateStructJenCode())
		if withPlural {
			structsCode = append(structsCode, s.GenerateSliceType())
		}
	}
	return structsCode
}
