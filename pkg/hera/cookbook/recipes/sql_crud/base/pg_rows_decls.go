package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/conditionals"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
)

// GeneratePgRowsPtrFn generates the templated PgRowsPtrFunc
func GeneratePgRowsPtrFn(structGen primitive.StructGen, body ...*jen.Statement) jen.Code {
	fnName := "GetRowValues"
	fn := funcs.NewFn(fnName)
	fn.Fields = genPgRowsFnFields()
	fn.ReturnFields = genPgRowsFnReturnFields()
	fn.AddBodyStatement(body...)
	return fn.GenerateStructPtrFunc(structGen)
}

// generateSwitchStatementForPgRows you will want to break this fn up when adding more than just the default case
// this generates the total switch, which adds all the case conditions
func generateSwitchStatementForPgRows(v vars.VariableGen, structGen primitive.StructGen) *jen.Statement {
	sc := conditionals.NewSwitchCase("queryName")
	// add other conditions here
	sc.AddCondition(generateDefaultCaseStatement(v, structGen))
	jc := sc.GenerateSwitchStatement()
	return jc
}

// generateDefaultCaseStatement is just getting all the fields as rows and using for the default
// create other fns to group complex cases
func generateDefaultCaseStatement(v vars.VariableGen, structGen primitive.StructGen) fields.CaseField {
	declDefaultCaseFields := declarePgValuesStructVar(v, "embedded", structGen)
	cf := fields.NewCaseField("default", "")
	cf.AddBodyStatement(declDefaultCaseFields)
	return cf
}
