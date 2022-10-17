package vars

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

type VariableGen struct {
	StringConstants map[string]string
	Structs         map[string]structs.StructGen
}

func NewVarGen() VariableGen {
	v := VariableGen{
		StringConstants: make(map[string]string),
		Structs:         make(map[string]structs.StructGen),
	}
	return v
}

func (v *VariableGen) CreateConstStringDecl(name string) *jen.Statement {
	if value, ok := v.StringConstants[name]; ok {
		return jen.Null().Const().Id(name).Op("=").Lit(value)
	}
	return jen.Null()
}

func (v *VariableGen) CreateStructDecl(varName, pkgName, structName string, isFirstInit bool) *jen.Statement {
	eq := "="
	if isFirstInit {
		eq = ":="
	}
	if _, ok := v.Structs[structName]; ok {
		return jen.Id(varName).Op(eq).Id(pkgName).Dot(structName).Values()
	}
	return jen.Null()
}

func (v *VariableGen) InsertStruct(s structs.StructGen) {
	if _, ok := v.Structs[s.Name]; !ok {
		v.Structs[s.Name] = structs.StructGen{}
	}
	v.Structs[s.Name] = s
}
