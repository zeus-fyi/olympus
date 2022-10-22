package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

func GetPgRowsTemplateDeclarations(structGen primitive.StructGen) (vars.VariableGen, primitive.StructGen, *jen.Statement) {
	v := genPgRowsVar()
	bodyInitPgRowsStruct := declarePgValuesStructVar(v, "init", v.GenStructInstructs["init"])
	return v, structGen, bodyInitPgRowsStruct
}

// declarePgValuesStructVar use the genKey to create more complex conditions and add their fields to the mapping
// this should be a more general fn later on and then uses the hardcoded values as constants in this pkg
func declarePgValuesStructVar(v vars.VariableGen, genKey string, embeddedStruct primitive.StructGen) *jen.Statement {
	declStruct := v.CreateEmbeddedStructVarDeclForSlice("pgValues", "apps", genKey, embeddedStruct)
	return declStruct
}

func genPgRowsVar() vars.VariableGen {
	v := vars.NewVarGen()
	v.InsertStruct(genRowValuesStructTemplate())
	return v
}

func genPgRowsFnFields() []fields.Field {
	return []fields.Field{{
		Name:  "queryName",
		Type:  "string",
		Value: "",
	}}
}

func genPgRowsFnReturnFields() []fields.Field {
	return []fields.Field{{
		Pkg:   "apps",
		Name:  "pgValues",
		Type:  "RowValues",
		Value: "",
	}}
}

func genRowValuesStructTemplate() primitive.StructGen {
	structToMake := primitive.StructGen{
		Name:   "RowValues",
		Fields: nil,
	}
	fieldOne := fields.Field{
		Name: "RowValues",
		Type: "struct",
	}
	structToMake.AddField(fieldOne)
	return structToMake
}
