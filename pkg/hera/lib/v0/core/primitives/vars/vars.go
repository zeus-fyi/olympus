package vars

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

type VariableGen struct {
	StringConstants      map[string]string
	StringSliceConstants map[string][]string

	// use the key to derive the gen logic
	GenStructInstructs map[string]structs.StructGen
}

func NewVarGen() VariableGen {
	v := VariableGen{
		StringConstants:      make(map[string]string),
		StringSliceConstants: make(map[string][]string),
		GenStructInstructs:   make(map[string]structs.StructGen),
	}
	return v
}

func (v *VariableGen) GenImportedVarReference(importedPkgName string) *jen.Statement {
	if varName, ok := v.StringConstants[importedPkgName]; ok {
		return jen.Id(importedPkgName).Dot(varName)
	}
	// TODO, handle local reference
	return jen.Null()
}

func (v *VariableGen) SetAndReturnImportedVarReference(importedPkgName, importedVarName string) VariableGen {
	if v == nil {
		// TODO should also set and return
		return NewVarGen()
	}
	v.StringConstants[importedPkgName] = importedVarName
	return *v
}

func (v *VariableGen) InsertStruct(s structs.StructGen) {
	if _, ok := v.GenStructInstructs[s.Name]; !ok {
		v.GenStructInstructs[s.Name] = structs.StructGen{}
	}
	v.GenStructInstructs[s.Name] = s
	v.GenStructInstructs["init"] = s
	v.GenStructInstructs["embedded"] = structs.StructGen{}
}
