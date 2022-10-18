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

func GenerateCaseStatementForPgRows(sg primitive.StructsGen) jen.Code {
	sc := conditionals.NewSwitchCase("queryName")

	//cf := fields.CaseField{
	//	Name: "default",
	//	Type: "default",
	//	Body: nil,
	//}
	//
	//for k, s := range sg.StructsMap {
	//
	//	sc.Conditions["default"].Body = s.GenerateStructJenStmt()
	//}
	jc := sc.GenerateSwitchStatement()
	return jc
}
