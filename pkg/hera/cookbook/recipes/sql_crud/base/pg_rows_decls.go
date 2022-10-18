package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/conditionals"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

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

func DeclarePgValuesStructVar(v vars.VariableGen, genKey string, embeddedStruct primitive.StructGen) *jen.Statement {
	declStruct := v.CreateEmbeddedStructVarDeclForSlice("pgValues", "apps", genKey, embeddedStruct)
	return declStruct
}

func generateDefaultCaseStatement(v vars.VariableGen, structGen primitive.StructGen) fields.CaseField {
	declDefaultCaseFields := DeclarePgValuesStructVar(v, "embedded", structGen)
	cf := fields.NewCaseField("default", "")
	cf.AddBodyStatement(declDefaultCaseFields)
	return cf
}

func GenerateSwitchStatementForPgRows(v vars.VariableGen, structGen primitive.StructGen) *jen.Statement {
	sc := conditionals.NewSwitchCase("queryName")
	sc.AddCondition(generateDefaultCaseStatement(v, structGen))
	jc := sc.GenerateSwitchStatement()
	return jc
}
