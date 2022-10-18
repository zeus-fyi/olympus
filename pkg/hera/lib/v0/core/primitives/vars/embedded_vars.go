package vars

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

func IsFirstInit(genKey string) string {
	eq := "="
	if genKey == "init" {
		eq = ":="
	}
	return eq
}

func (v *VariableGen) CreateEmbeddedStructVarDeclForSlice(varName, pkgName, genKey string, embeddedStruct structs.StructGen) *jen.Statement {
	eq := IsFirstInit(genKey)
	ws, ok := v.GenStructInstructs["init"]
	if !ok {
		return jen.Null()
	}
	var jenSlice []jen.Code
	jenCode := jen.Id(varName).Op(eq).Id(pkgName).Dot(ws.Name)
	switch genKey {
	case "init":
		return jenCode.Values()
	case "embedded":
		esh := embeddedStruct.ShortHand()
		for _, f := range embeddedStruct.Fields {
			jenSlice = append(jenSlice, jen.Id(esh).Dot(f.Name))
		}
	default:
	}
	return jen.Id(varName).Op(eq).Id(pkgName).Dot(ws.Name).Values(jenSlice...)
}
